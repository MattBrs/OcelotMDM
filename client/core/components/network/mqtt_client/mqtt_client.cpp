#include "mqtt_client.hpp"

#include <mqtt/async_client.h>
#include <mqtt/buffer_ref.h>

#include <algorithm>
#include <format>
#include <string>

namespace OcelotMDM::component::network {
MqttClient::MqttClient(
    const std::string &host, const std::uint32_t port,
    const std::string &clientID)
    : client(std::format("tcp://{}:{}", host, std::to_string(port)), clientID) {
    this->connectOpts.set_clean_session(false);
    // this->connectOpts.set_automatic_reconnect(true);

    this->client.set_disconnected_handler(
        [this](const mqtt::properties &props, int reasonCode) {
            this->reconnect();
        });
}

bool MqttClient::connect() {
    auto connRes =
        this->client.connect(this->connectOpts)->wait_for(MQTT_TIMEOUT);
    if (!connRes) {
        // connection timed out
        return false;
    }

    return true;
}

bool MqttClient::subscribe(const std::string &topic, const std::uint32_t qos) {
    auto subToken = this->client.subscribe(topic, qos);
    auto waitRes = subToken->wait_for(MQTT_TIMEOUT);

    if (!waitRes) {
        return false;
    }

    if (!this->topics.contains(topic)) {
        this->topics[topic] = true;
    }

    return true;
}

bool MqttClient::publish(
    const std::string &msg, const std::string &topic, const std::uint32_t qos) {
    auto pubToken =
        this->client.publish(topic, msg.data(), msg.size(), qos, true);
    auto waitRes = pubToken->wait_for(MQTT_TIMEOUT);

    if (!waitRes) {
        return false;
    }

    return true;
}

bool MqttClient::disconnect() {
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

void MqttClient::reconnect() {
    auto res = this->connect();
    if (!res) {
        return;
    }

    // sub to each topic again
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
}  // namespace OcelotMDM::component::network
