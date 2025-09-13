#pragma once

#include <atomic>
#include <condition_variable>
#include <cstdint>
#include <mutex>
#include <optional>
#include <string>
#include <thread>

#include "mqtt_client.hpp"

namespace OcelotMDM::component::service {
class UptimeService {
   public:
    UptimeService(
        const std::string &mqttIp, const std::uint32_t port,
        const std::string &deviceID);

    ~UptimeService();
    void workerFunction();

   private:
    const int UPDATE_FREQUENCY = 60000 * 2;

    std::string         deviceID;
    network::MqttClient mqttClient;

    std::thread             workerTh;
    std::condition_variable workerCv;
    std::mutex              workerMtx;
    std::atomic<bool>       shouldStop;

    std::optional<std::string> getVpnAddress();
};
};  // namespace OcelotMDM::component::service
