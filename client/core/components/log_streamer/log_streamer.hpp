#pragma once

#include <atomic>
#include <condition_variable>
#include <memory>
#include <mutex>
#include <queue>
#include <string>
#include <thread>

#include "mqtt_client.hpp"

namespace OcelotMDM::component {
class LogStreamer {
   public:
    LogStreamer(
        const std::shared_ptr<network::MqttClient> &mqttClient,
        const std::string                          &deviceID);
    ~LogStreamer();

    std::shared_ptr<std::queue<std::string>> getQueue();

    void run();
    void stop();

    bool isRunning();

   private:
    const int WORKER_QUEUE_INTR_MS = 500;

    std::string deviceID;

    std::shared_ptr<std::queue<std::string>> logQueue;
    std::shared_ptr<network::MqttClient>     mqttClient;

    std::thread             wrkTh;
    std::condition_variable wrkCv;
    std::mutex              wrkMtx;
    std::atomic<bool>       wrkRunning;

    void queueWorker();
};
}  // namespace OcelotMDM::component
