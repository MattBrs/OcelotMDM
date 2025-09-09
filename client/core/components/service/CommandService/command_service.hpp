#pragma once

#include <mqtt/message.h>

#include <cstdint>
#include <string>

#include "command_model.hpp"
#include "mqtt_client.hpp"

namespace OcelotMDM::component::service {
class CommandService {
   public:
    CommandService(
        const std::string &mqttIp, const std::uint32_t port,
        const std::string &deviceID);

   private:
    network::MqttClient mqttClient;

    model::Command decodeCmdMsg(mqtt::const_message_ptr msg);

    void onCmdArrived(mqtt::const_message_ptr);
};
};  // namespace OcelotMDM::component::service
