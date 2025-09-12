#include "command_service.hpp"

#include <mqtt/message.h>

#include <cstdint>
#include <functional>
#include <iomanip>
#include <ios>
#include <list>
#include <memory>
#include <nlohmann/json_fwd.hpp>
#include <ostream>
#include <sstream>
#include <string>
#include <vector>

#include "command_dao.hpp"
#include "command_model.hpp"
#include "mqtt_client.hpp"
#include "nlohmann/json.hpp"

namespace OcelotMDM::component::service {
CommandService::CommandService(
    const std::shared_ptr<db::CommandDao> &cmdDao, const std::string &mqttIp,
    const std::uint32_t mqttPort, const std::string &deviceID)
    : cmdDao(cmdDao),
      deviceID(deviceID),
      mqttClient(mqttIp, mqttPort, deviceID) {
    auto queuedCommands = this->cmdDao->getQueuedCommands();
    if (queuedCommands.has_value()) {
        for (auto &cmd : queuedCommands.value()) {
            std::cout << "queueing command from db: " << cmd.id << std::endl;
            this->cmdQueue.push(cmd);
            this->queuedCmds.emplace(cmd.id);
        }
    }

    this->mqttClient.setMsgCallback(
        std::bind(&CommandService::onCmdArrived, this, std::placeholders::_1));

    this->mqttClient.connect();
    mqttClient.subscribe(deviceID + "/cmd", 1);
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

    std::cout << "arrived from mqtt: " << rawData.dump() << std::endl;
    model::Command cmd{
        rawData["Id"], rawData["MessageAction"], rawData["Payload"],
        rawData["Priority"]};

    return cmd;
}

void CommandService::onCmdArrived(mqtt::const_message_ptr msg) {
    auto cmd = this->decodeCmdMsg(msg);
    std::cout << "arrived cmd: \n";
    std::cout << "- " << cmd.id << "\n";
    std::cout << "- " << cmd.commandAction << "\n";
    std::cout << "- " << cmd.payload << "\n";
    std::cout << "- " << cmd.priority << "\n" << std::flush;

    this->enqueueCommand(cmd);
}

void CommandService::enqueueCommand(const model::Command &cmd) {
    if (this->queuedCmds.contains(cmd.id)) {
        // command is already queued
        std::cout << "cmd arrived but already in queue, skipping: " << cmd.id
                  << std::endl;
        return;
    }

    this->cmdQueue.push(cmd);
    this->queuedCmds.emplace(cmd.id);

    auto ackRes = nlohmann::json();
    ackRes["Id"] = cmd.id;
    ackRes["State"] = "acknowledged";
    ackRes["errorMsg"] = "";

    auto encoded = bytesToHex(nlohmann::json::to_msgpack(ackRes));

    std::cout << "about to pub" << std::endl;
    auto pubRes = this->mqttClient.publish(encoded, this->deviceID + "/ack", 1);

    this->cmdDao->enqueueCommand(cmd);
}
}  // namespace OcelotMDM::component::service
