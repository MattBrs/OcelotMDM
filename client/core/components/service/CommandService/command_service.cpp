#include "command_service.hpp"

#include <mqtt/message.h>

#include <cstdint>
#include <functional>
#include <list>
#include <nlohmann/json_fwd.hpp>
#include <ostream>
#include <string>
#include <vector>

#include "command_model.hpp"
#include "mqtt_client.hpp"
#include "nlohmann/json.hpp"

namespace OcelotMDM::component::service {
CommandService::CommandService(
    const std::string &mqttIp, const std::uint32_t mqttPort,
    const std::string &deviceID)
    : mqttClient(mqttIp, mqttPort, deviceID) {
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

model::Command CommandService::decodeCmdMsg(mqtt::const_message_ptr msg) {
    auto rawData = nlohmann::json::from_msgpack(hexToBytes(msg->to_string()));

    std::cout << "arrived from mqtt: " << rawData.dump() << std::endl;
    model::Command cmd{
        .id = rawData["Id"],
        .commandAction = rawData["MessageAction"],
        .payload = rawData["Payload"],
        .priority = rawData["Priority"]};

    return cmd;
}

void CommandService::onCmdArrived(mqtt::const_message_ptr msg) {
    auto cmd = this->decodeCmdMsg(msg);
    std::cout << "arrived cmd: \n";
    std::cout << "- " << cmd.id << "\n";
    std::cout << "- " << cmd.commandAction << "\n";
    std::cout << "- " << cmd.payload << "\n";
    std::cout << "- " << cmd.priority << "\n" << std::flush;

    // todo: send ack to backend
}
}  // namespace OcelotMDM::component::service
