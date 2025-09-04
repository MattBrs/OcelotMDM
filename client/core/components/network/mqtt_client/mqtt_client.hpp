#pragma once

#include <mqtt/async_client.h>
#include <mqtt/connect_options.h>

#include <cstdint>
#include <string>
#include <unordered_map>
#include <vector>

namespace OcelotMDM::component::network {
// synced implementation of paho mqtt client
class MqttClient {
   public:
    MqttClient(
        const std::string &host, const std::uint32_t port,
        const std::string &clientID, const std::vector<std::string> &topics);
    ~MqttClient();

    bool connect();

    bool subscribe(const std::string &topic, const std::uint32_t qos);

    bool publish(
        const std::string &msg, const std::string &topic,
        const std::uint32_t qos);

    bool disconnect();

   private:
    const int MQTT_TIMEOUT = 10000;

    std::string   host;
    std::uint32_t port;
    std::string   clientID;

    mqtt::async_client    client;
    mqtt::connect_options connectOpts;

    std::unordered_map<std::string, bool> topics;

    void subscribeTopics();
};
};  // namespace OcelotMDM::component::network
