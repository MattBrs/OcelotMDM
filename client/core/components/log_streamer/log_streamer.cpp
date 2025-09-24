#include "log_streamer.hpp"

#include <chrono>
#include <memory>
#include <mutex>
#include <string>
#include <thread>

namespace OcelotMDM::component {
LogStreamer::LogStreamer(
    const std::shared_ptr<network::MqttClient> &mqttClient,
    const std::string                          &deviceID)
    : deviceID(deviceID) {
    this->mqttClient = mqttClient;
    this->logQueue = std::make_shared<std::queue<std::string>>();
}

LogStreamer::~LogStreamer() {
    this->stop();
}

std::shared_ptr<std::queue<std::string>> LogStreamer::getQueue() {
    return this->logQueue;
}

void LogStreamer::run() {
    if (this->wrkRunning.load()) {
        return;
    }

    this->mqttClient->publish(
        "----start live logs---- ",
        std::string{this->deviceID}.append("/live-logs"), 1);

    this->wrkRunning.store(true);
    this->wrkTh = std::thread(&LogStreamer::queueWorker, this);
}

void LogStreamer::stop() {
    this->mqttClient->publish(
        "----stop live logs---- ",
        std::string{this->deviceID}.append("/live-logs"), 1);

    this->wrkRunning.store(false);
    this->wrkCv.notify_one();

    if (this->wrkTh.joinable()) {
        this->wrkTh.join();
    }
}

bool LogStreamer::isRunning() {
    return this->wrkRunning.load();
}

void LogStreamer::queueWorker() {
    while (this->wrkRunning.load()) {
        if (!this->logQueue->empty()) {
            auto item = this->logQueue->front();
            this->logQueue->pop();

            this->mqttClient->publish(
                item, std::string{this->deviceID}.append("/live-logs"), 1);
        }

        std::unique_lock<std::mutex> lock(this->wrkMtx);
        this->wrkCv.wait_for(
            lock, std::chrono::milliseconds(this->WORKER_QUEUE_INTR_MS),
            [this]() { return !this->wrkRunning.load(); });
    }
}
}  // namespace OcelotMDM::component
