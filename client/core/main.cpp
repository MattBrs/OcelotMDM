#include <iostream>

#include "http_client.hpp"
#include "mqtt_client.hpp"

using HttpClient = OcelotMDM::component::network::HttpClient;
using MqttClient = OcelotMDM::component::network::MqttClient;

int main() {
    MqttClient mqttClient("159.89.2.75", 1883, "sugo-boy", {"sugo-boy/cmd"});

    std::thread mqttTh([&mqttClient]() { mqttClient.connect(); });
    mqttTh.detach();

    HttpClient httpClient("https://httpbin.org/ip");

    std::cout << "Hello!" << std::endl;

    while (true);
}
