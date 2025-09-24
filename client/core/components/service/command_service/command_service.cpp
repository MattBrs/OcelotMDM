#include "command_service.hpp"

#include <mqtt/message.h>

#include <chrono>
#include <cstdint>
#include <functional>
#include <iomanip>
#include <ios>
#include <memory>
#include <mutex>
#include <nlohmann/json_fwd.hpp>
#include <optional>
#include <ostream>
#include <sstream>
#include <string>
#include <thread>
#include <vector>

#include "command_dao.hpp"
#include "command_model.hpp"
#include "commands_impl.hpp"
#include "http_client.hpp"
#include "log_streamer.hpp"
#include "logger.hpp"
#include "mqtt_client.hpp"
#include "nlohmann/json.hpp"

namespace OcelotMDM::component::service {
CommandService::CommandService(
    const std::shared_ptr<db::CommandDao>      &cmdDao,
    const std::shared_ptr<network::MqttClient> &mqttClient,
    const std::string &httpBaseUrl, const std::string &deviceID)
    : cmdDao(cmdDao),
      mqttClient(mqttClient),
      deviceID(deviceID),
      httpClient(httpBaseUrl),
      logStreamer(mqttClient, deviceID) {
    auto queuedCommands = this->cmdDao->getQueuedCommands();
    if (queuedCommands.has_value()) {
        for (auto &cmd : queuedCommands.value()) {
            Logger::getInstance().put(
                "queueing command from db: " + cmd.getId());
            this->cmdQueue.push(cmd);
            this->queuedCmds.emplace(cmd.getId());
        }
    }

    this->mqttClient->setMsgCallback(
        std::bind(&CommandService::onCmdArrived, this, std::placeholders::_1));

    mqttClient->subscribe(deviceID + "/cmd", 1);

    this->shouldStopTh.store(false);
    this->queueTh = std::thread(&CommandService::queueWorker, this);
}

CommandService::~CommandService() {
    this->shouldStopTh.store(true);
    this->queueCv.notify_one();

    if (this->queueTh.joinable()) {
        this->queueTh.join();
    }
}

void CommandService::queueWorker() {
    while (!this->shouldStopTh.load()) {
        if (this->cmdQueue.size() > 0) {
            auto cmd = this->cmdQueue.top();
            this->cmdQueue.pop();
            this->queuedCmds.erase(cmd.getId());
            auto execRes = this->executeCommand(cmd);

            if (!execRes.has_value()) {
                cmd.setError("command not supported");
                cmd.setStatus(model::Command::CommandStatus::Errored);
            } else if (!execRes.value().successful) {
                cmd.setError(execRes.value().props.error);
                cmd.setStatus(model::Command::CommandStatus::Errored);
            } else {
                cmd.setError("");
                cmd.setStatus(model::Command::CommandStatus::Completed);
            }

            auto encoded = this->encodeCmd(cmd);
            this->mqttClient->publish(encoded, this->deviceID + "/ack", 1);

            this->cmdDao->dequeCommand(cmd.getId());
        }

        std::unique_lock<std::mutex> lock(this->queueMtx);
        this->queueCv.wait_for(
            lock, std::chrono::milliseconds(DEQUEUE_INTR),
            [this]() { return this->shouldStopTh.load(); });
    }
}

std::vector<std::uint8_t> hexToBytes(const std::string &hex) {
    std::vector<std::uint8_t> bytes;

    if (hex.size() % 2 != 0) {
        // odd sized hex. not possible
        return bytes;
    }

    for (auto i = 0; i < hex.size(); i += 2) {
        bytes.emplace_back(std::stoi(hex.substr(i, 2), nullptr, 16));
    }

    return bytes;
}

std::string bytesToHex(const std::vector<std::uint8_t> &bytes) {
    std::stringstream ss;
    auto              start = bytes.begin();
    auto              end = bytes.end();

    ss << std::hex << std::setw(2) << std::setfill('0');

    while (start != end) {
        ss << static_cast<unsigned>(*start++);
    }

    return ss.str();
}

model::Command CommandService::decodeCmdMsg(mqtt::const_message_ptr msg) {
    auto rawData = nlohmann::json::from_msgpack(hexToBytes(msg->to_string()));

    model::Command cmd{
        rawData["Id"], rawData["MessageAction"], rawData["Payload"],
        rawData["Priority"], rawData["ReqiuredOnline"]};

    return cmd;
}

void CommandService::onCmdArrived(mqtt::const_message_ptr msg) {
    auto cmd = this->decodeCmdMsg(msg);

    this->enqueueCommand(cmd);
}

void CommandService::enqueueCommand(model::Command &cmd) {
    if (this->queuedCmds.contains(cmd.getId())) {
        return;
    }

    auto insertRes = this->cmdDao->enqueueCommand(cmd);
    if (!insertRes.has_value()) {
        return;
    }

    if (insertRes.value()) {
        this->cmdQueue.push(cmd);
        this->queuedCmds.emplace(cmd.getId());
    } else {
        Logger::getInstance().putError("errored on command insert");
        cmd.setStatus(model::Command::CommandStatus::Errored);
        cmd.setError(this->cmdDao->getError().c_str());
    }

    auto encoded = this->encodeCmd(cmd);
    this->mqttClient->publish(encoded, this->deviceID + "/ack", 1);
}

std::string CommandService::encodeCmd(const model::Command &cmd) {
    auto ackRes = nlohmann::json();
    ackRes["Id"] = cmd.getId();
    ackRes["State"] = cmd.getStatus();
    ackRes["ErrorMsg"] = cmd.getError();

    return bytesToHex(nlohmann::json::to_msgpack(ackRes));
}

std::optional<CommandImpl::ExecutionResult> CommandService::executeCommand(
    const model::Command &cmd) {
    Logger::getInstance().put("command action " + cmd.getAction());

    if (cmd.getAction().compare("install_binary") == 0) {
        nlohmann::json payload;
        try {
            payload = nlohmann::json::parse(cmd.getPayload());
        } catch (nlohmann::json::exception &e) {
            return std::nullopt;
        }

        auto res = CommandImpl::installBinary(
            &this->httpClient, payload["name"], payload["otp"]);
        return res;
    }

    if (cmd.getAction().compare("send_logs") == 0) {
        auto res = CommandImpl::sendLogs(this->mqttClient, this->deviceID);
        return res;
    }

    if (cmd.getAction().compare("enable_live_logging") == 0) {
        CommandImpl::ExecutionResult res;
        if (this->logStreamer.isRunning()) {
            res.successful = false;
            res.props.error = "live logging already enabled";

            return res;
        }

        auto streamerQueue = this->logStreamer.getQueue();
        Logger::getInstance().registerQueue(streamerQueue);
        this->logStreamer.run();

        this->timer.start(
            [this]() {
                Logger::getInstance().registerQueue(nullptr);
                this->logStreamer.stop();
            },
            1 * 60 * 1000, false);

        res.successful = true;
        return res;
    }

    return std::nullopt;
}
}  // namespace OcelotMDM::component::service
