#include "command_service.hpp"

#include <mqtt/message.h>

#include <cstdint>
#include <list>
#include <nlohmann/json_fwd.hpp>
#include <string>
#include <vector>

#include "mqtt_client.hpp"
#include "nlohmann/json.hpp"

namespace OcelotMDM::component::service {
CommandService::CommandService(
    const std::string &mqttIp, const std::uint32_t mqttPort,
    const std::string &deviceID)
    : mqttClient(mqttIp, mqttPort, deviceID) {
    mqttClient.subscribe(deviceID + "/cmd", 1);

    this->mqttClient.setMsgCallback([this](mqtt::const_message_ptr msg) {
        // decode msg and add to queue
        // todo: define msg structure

        this->decodeCmdMsg(msg);
    });
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

void CommandService::decodeCmdMsg(mqtt::const_message_ptr msg) {
    auto rawData = nlohmann::json::from_msgpack(hexToBytes(msg->to_string()));

    std::cout << "arrived from mqtt: " << rawData.dump() << std::endl;
    /*
      * msgContent:
         Id.Hex(),
         CommandActionName,
         Payload,
         Priority,
      */
}
}  // namespace OcelotMDM::component::service
