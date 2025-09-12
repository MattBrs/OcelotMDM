#include <iostream>

#include "command_service.hpp"
#include "db_client.hpp"
#include "http_client.hpp"
#include "mqtt_client.hpp"

using HttpClient = OcelotMDM::component::network::HttpClient;
using MqttClient = OcelotMDM::component::network::MqttClient;
using CommandService = OcelotMDM::component::service::CommandService;
using DbClient = OcelotMDM::component::db::DbClient;

int main() {
    std::cout << "Starting!" << std::endl;

    DbClient dbClient("test.db");

    CommandService cmdService(
        dbClient.getCommandDao(), "159.89.2.75", 1883, "misty-dew");
    HttpClient httpClient("https://httpbin.org/ip");

    std::cout << "Hello!" << std::endl;

    while (true);
}
