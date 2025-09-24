#include "mqtt_client.hpp"

#include <mqtt/async_client.h>
#include <mqtt/buffer_ref.h>
#include <mqtt/exception.h>
#include <mqtt/message.h>

#include <algorithm>
#include <chrono>
#include <format>
#include <iostream>
#include <mutex>
#include <string>
#include <utility>

#include "logger.hpp"

namespace OcelotMDM::component::network {
MqttClient::MqttClient(
    const std::string &host, const std::uint32_t port,
    const std::string &clientID, const std::vector<std::string> &topics)
    : client(std::format("tcp://{}:{}", host, std::to_string(port)), clientID) {
    this->connectOpts.set_clean_session(false);
    this->connectOpts.set_keep_alive_interval(2);

    this->client.set_connected_handler([this](const std::basic_string<char> &) {
        Logger::getInstance().put("mqtt client connected successfully");

        this->connected.store(true);
        this->queueCv.notify_one();
    });

    this->client.set_disconnected_handler(
        [this](const mqtt::properties &props, int reasonCode) {
            Logger::getInstance().put("mqtt client disconnected");

            this->connected.store(false);
            // this->reconnectCv.notify_one();
        });

    this->client.set_connection_lost_handler(
        [this](const std::basic_string<char> &) {
            Logger::getInstance().put("mqtt client lost connection");

            this->connected.store(false);
            // this->reconnectCv.notify_one();
        });

    for (auto &item : topics) {
        this->topics[item] = false;
    }

    this->client.set_message_callback([this](mqtt::const_message_ptr msg) {
        if (this->msgArrivedCb != nullptr) {
            this->msgArrivedCb(msg);
            return;
        }

        std::cout << "message arrived on topic: " << msg->get_topic()
                  << " but no one listened to it\n";
    });

    this->reconnectTh = std::thread(&MqttClient::reconnectionWorker, this);
    this->queueTh = std::thread(&MqttClient::queueWorker, this);
}

bool MqttClient::connect() {
    Logger::getInstance().put("trying connection to mqtt broker");

    if (this->client.is_connected()) {
        Logger::getInstance().put("already connected");
        this->connected.store(true);
        return true;
    }

    try {
        auto connRes =
            this->client.connect(this->connectOpts)->wait_for(MQTT_TIMEOUT);
        if (!connRes) {
            this->connected.store(false);
            Logger::getInstance().putError("connection timed out");
            return false;
        }
    } catch (const mqtt::exception &e) {
        this->connected.store(false);
        Logger::getInstance().putError(
            "connection failed: " + std::string{e.what()});
        return false;
    }

    this->connected.store(true);
    this->subscribeTopics();
    return true;
}

bool MqttClient::subscribe(const std::string &topic, const std::uint32_t qos) {
    auto guard = std::lock_guard<std::mutex>(this->usageMtx);

    if (!this->connected.load()) {
        this->topics[topic] = false;
        return false;
    }

    auto subToken = this->client.subscribe(topic, qos);
    auto waitRes = subToken->wait_for(MQTT_TIMEOUT);

    this->topics[topic] = waitRes;
    return waitRes;
}

bool MqttClient::publish(
    const std::string &msg, const std::string &topic, const std::uint32_t qos,
    const bool retain) {
    auto guard = std::lock_guard<std::mutex>(this->usageMtx);

    if (!this->connected.load()) {
        if (!retain) {
            // we dont care about retained messages
            // because those are for current device status.
            Logger::getInstance().put("tried publish, was offline. enqueued");
            this->messageQueue.emplace(std::make_pair(topic, msg));
        }

        return false;
    }

    auto pubToken =
        this->client.publish(topic, msg.data(), msg.size(), qos, retain);
    pubToken->try_wait();

    return true;
}

bool MqttClient::disconnect() {
    if (!this->client.is_connected()) {
        return true;
    }

    auto discRes = this->client.disconnect()->wait_for(MQTT_TIMEOUT);
    if (!discRes) {
        // could not disconnect
        return false;
    }

    // set topics as not subbed to
    std::for_each(
        this->topics.begin(), this->topics.end(),
        [](std::pair<const std::string, bool> &item) { item.second = false; });

    return true;
}

void MqttClient::subscribeTopics() {
    if (!this->client.is_connected()) {
        return;
    }

    std::for_each(
        this->topics.begin(), this->topics.end(),
        [this](std::pair<const std::string, bool> &item) {
            this->subscribe(item.first, 1);
        });
}

MqttClient::~MqttClient() {
    this->shouldStopTh.store(true);
    this->reconnectCv.notify_one();

    if (this->reconnectTh.joinable()) {
        this->reconnectTh.join();
    }

    this->client.set_disconnected_handler(
        [this](const mqtt::properties &props, int reasonCode) {
            // do nothing
        });

    this->disconnect();
}

void MqttClient::setMsgCallback(
    std::function<void(mqtt::const_message_ptr)> cb) {
    this->msgArrivedCb = cb;
}

void MqttClient::reconnectionWorker() {
    while (!this->shouldStopTh.load()) {
        std::unique_lock<std::mutex> lock(this->reconnectMtx);
        this->reconnectCv.wait_for(
            lock, std::chrono::milliseconds(this->MQTT_CONN_CHECK),
            [this]() { return this->shouldStopTh.load(); });

        if (this->shouldStopTh.load()) {
            continue;
        }

        Logger::getInstance().put(
            "woke up and im going to figure out if we are connected");

        if (this->connected.load()) {
            Logger::getInstance().put("still connected");
            continue;
        }

        Logger::getInstance().put("oopsie, we are disconnected");

        this->connect();
    }
}

void MqttClient::queueWorker() {
    while (!this->shouldStopTh.load()) {
        while (!this->messageQueue.empty() && this->connected.load() &&
               !this->shouldStopTh.load()) {
            Logger::getInstance().put("dequeing a message");
            auto msg = this->messageQueue.front();
            this->messageQueue.pop();

            this->publish(msg.second, msg.first, 1);
        }

        std::unique_lock<std::mutex> lock(this->queueMtx);
        this->queueCv.wait(lock, [this]() {
            return this->shouldStopTh.load() || !this->messageQueue.empty();
        });
    }
}
}  // namespace OcelotMDM::component::network
