#include <cstdint>
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <memory>
#include <nlohmann/json.hpp>
#include <nlohmann/json_fwd.hpp>

#include "command_service.hpp"
#include "components/service/uptime_service/uptime_service.hpp"
#include "db_client.hpp"
#include "http_client.hpp"
#include "logger.hpp"
#include "mqtt_client.hpp"
#include "utils.hpp"

using Logger = OcelotMDM::component::Logger;
using DbClient = OcelotMDM::component::db::DbClient;
using HttpClient = OcelotMDM::component::network::HttpClient;
using MqttClient = OcelotMDM::component::network::MqttClient;
using CommandService = OcelotMDM::component::service::CommandService;
using UptimeService = OcelotMDM::component::service::UptimeService;

int main() {
    std::cout << "Starting!" << std::endl;

    auto conf = nlohmann::json();
    try {
        auto fileContent = OcelotMDM::component::utils::readFile("conf.json");
        if (!fileContent.has_value()) {
            Logger::getInstance().put("could not find conf file");
            exit(1);
        }

        conf = nlohmann::json::parse(fileContent.value());
    } catch (const nlohmann::json::exception &e) {
        Logger::getInstance().put("could not read conf file");
        exit(1);
    }

    const std::uint32_t mqttPort = 1883;
    const std::string   deviceName = conf["device_name"];
    const std::string   apiBaseUrl = conf["api_base_url"];
    const std::string   mqttHost = conf["mqtt_host"];

    auto mqttClient =
        std::make_shared<MqttClient>(mqttHost, mqttPort, deviceName);
    mqttClient->connect();

    DbClient       dbClient("test.db");
    UptimeService  uptimeService(mqttClient, deviceName);
    CommandService cmdService(
        dbClient.getCommandDao(), mqttClient, apiBaseUrl, deviceName);

    while (true);
}
