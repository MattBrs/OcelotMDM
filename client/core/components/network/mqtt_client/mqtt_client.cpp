#include "mqtt_client.hpp"

#include <mqtt/async_client.h>
#include <mqtt/buffer_ref.h>
#include <mqtt/message.h>

#include <algorithm>
#include <format>
#include <iostream>
#include <string>

namespace OcelotMDM::component::network {
MqttClient::MqttClient(
    const std::string &host, const std::uint32_t port,
    const std::string &clientID, const std::vector<std::string> &topics)
    : client(std::format("tcp://{}:{}", host, std::to_string(port)), clientID) {
    this->connectOpts.set_clean_session(false);
    // this->connectOpts.set_automatic_reconnect(true);

    this->client.set_disconnected_handler(
        [this](const mqtt::properties &props, int reasonCode) {
            this->connect();
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
}

bool MqttClient::connect() {
    if (this->client.is_connected()) {
        return true;
    }

    auto connRes =
        this->client.connect(this->connectOpts)->wait_for(MQTT_TIMEOUT);
    if (!connRes) {
        // connection timed out
        return false;
    }

    this->subscribeTopics();
    return true;
}

bool MqttClient::subscribe(const std::string &topic, const std::uint32_t qos) {
    if (!this->client.is_connected()) {
        this->topics[topic] = false;
        return false;
    }

    auto subToken = this->client.subscribe(topic, qos);
    auto waitRes = subToken->wait_for(MQTT_TIMEOUT);

    this->topics[topic] = waitRes;
    return waitRes;
}

bool MqttClient::publish(
    const std::string &msg, const std::string &topic, const std::uint32_t qos) {
    if (!this->client.is_connected()) {
        return false;
    }

    auto pubToken =
        this->client.publish(topic, msg.data(), msg.size(), qos, true);
    auto waitRes = pubToken->wait_for(MQTT_TIMEOUT);

    if (!waitRes) {
        return false;
    }

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
}  // namespace OcelotMDM::component::network
