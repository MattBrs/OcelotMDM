#include "command_service.hpp"

#include <mqtt/message.h>

#include <chrono>
#include <cstdint>
#include <functional>
#include <iomanip>
#include <ios>
#include <iostream>
#include <memory>
#include <mutex>
#include <nlohmann/json_fwd.hpp>
#include <optional>
#include <ostream>
#include <sstream>
#include <string>
#include <thread>
#include <vector>

#include "binary_dao.hpp"
#include "command_dao.hpp"
#include "command_model.hpp"
#include "commands_impl.hpp"
#include "http_client.hpp"
#include "linux_spawner_service.hpp"
#include "log_streamer.hpp"
#include "logger.hpp"
#include "mqtt_client.hpp"
#include "nlohmann/json.hpp"

namespace OcelotMDM::component::service {
CommandService::CommandService(
    const std::shared_ptr<db::CommandDao>      &cmdDao,
    const std::shared_ptr<db::BinaryDao>       &binDao,
    const std::shared_ptr<network::MqttClient> &mqttClient,
    const std::string &httpBaseUrl, const std::string &deviceID)
    : timer(std::make_shared<Timer>()),
      cmdDao(cmdDao),
      binDao(binDao),
      mqttClient(mqttClient),
      logStreamer(std::make_shared<LogStreamer>(mqttClient, deviceID)),
      deviceID(deviceID),
      httpClient(std::make_shared<network::HttpClient>(httpBaseUrl)) {
    // if linux use LinuxSpawner, else use AndroidSpawner
    this->spawnerService = std::make_unique<LinuxSpawerService>(this->binDao);

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

        CommandImpl::ExecutionResult res;
        try {
            res = CommandImpl::installBinary(
                this->binDao, this->httpClient, payload["name"],
                payload["otp"]);

            this->spawnerService->runBinary(res.props.applicationPath);
        } catch (nlohmann::json::exception &e) {
            res.successful = false;
            res.props.error = "could not parse payload";
        }

        return res;
    }

    if (cmd.getAction().compare("send_logs") == 0) {
        return CommandImpl::sendLogs(this->mqttClient, this->deviceID);
    }

    if (cmd.getAction().compare("enable_live_logging") == 0) {
        CommandImpl::ExecutionResult res;
        if (this->logStreamer == nullptr || this->timer == nullptr) {
            res.successful = false;
            res.props.error = "streamer or timer are not initialized";
            return res;
        }

        res = CommandImpl::enableLiveLogging(this->logStreamer, this->timer);
        return res;
    }

    if (cmd.getAction().compare("disable_live_logging") == 0) {
        CommandImpl::ExecutionResult res;
        if (this->logStreamer == nullptr || this->timer == nullptr) {
            res.successful = false;
            res.props.error = "streamer or timer are not initialized";
            return res;
        }

        res = CommandImpl::disableLiveLogging(this->logStreamer, this->timer);
        return res;
    }

    if (cmd.getAction().compare("start_binary") == 0) {
    }

    if (cmd.getAction().compare("uninstall_binary") == 0) {
    }

    if (cmd.getAction().compare("list_binaries") == 0) {
        return CommandImpl::listBinaries(this->binDao);
    }

    return std::nullopt;
}
}  // namespace OcelotMDM::component::service
