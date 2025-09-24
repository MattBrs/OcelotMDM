#pragma once

#include <mqtt/async_client.h>
#include <mqtt/connect_options.h>
#include <mqtt/message.h>

#include <atomic>
#include <condition_variable>
#include <cstdint>
#include <functional>
#include <mutex>
#include <queue>
#include <string>
#include <thread>
#include <unordered_map>
#include <vector>

namespace OcelotMDM::component::network {
// synced implementation of paho mqtt client
class MqttClient {
   public:
    MqttClient(
        const std::string &host, const std::uint32_t port,
        const std::string              &clientID,
        const std::vector<std::string> &topics = {});

    ~MqttClient();

    bool connect();

    bool subscribe(const std::string &topic, const std::uint32_t qos);

    bool publish(
        const std::string &msg, const std::string &topic,
        const std::uint32_t qos, const bool retain = false);

    bool disconnect();

    void setMsgCallback(std::function<void(mqtt::const_message_ptr)>);

   private:
    const int MQTT_TIMEOUT = 10000;
    const int MQTT_CONN_CHECK = 10000;

    std::string       host;
    std::uint32_t     port;
    std::string       clientID;
    std::mutex        usageMtx;
    std::atomic<bool> connected = false;

    std::condition_variable reconnectCv;
    std::thread             reconnectTh;
    std::mutex              reconnectMtx;
    std::atomic<bool>       shouldStopTh;

    std::condition_variable queueCv;
    std::thread             queueTh;
    std::mutex              queueMtx;

    mqtt::async_client    client;
    mqtt::connect_options connectOpts;

    std::queue<std::pair<std::string, std::string>> messageQueue;

    std::unordered_map<std::string, bool>            topics;
    std::function<void(mqtt::const_message_ptr msg)> msgArrivedCb = nullptr;

    void subscribeTopics();
    void reconnectionWorker();
    void queueWorker();
};
};  // namespace OcelotMDM::component::network
