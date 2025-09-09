#pragma once

#include <mqtt/message.h>

#include <cstdint>
#include <list>
#include <string>

#include "mqtt_client.hpp"

namespace OcelotMDM::component::service {
class CommandService {
   public:
    CommandService(
        const std::string &mqttIp, const std::uint32_t port,
        const std::string &deviceID);

   private:
    network::MqttClient mqttClient;

    void decodeCmdMsg(mqtt::const_message_ptr msg);
};
};  // namespace OcelotMDM::component::service
