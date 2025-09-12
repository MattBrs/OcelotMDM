#pragma once

#include <mqtt/message.h>

#include <cstdint>
#include <memory>
#include <queue>
#include <set>
#include <string>

#include "command_dao.hpp"
#include "command_model.hpp"
#include "mqtt_client.hpp"

namespace OcelotMDM::component::service {
class CommandService {
   public:
    CommandService(
        const std::shared_ptr<db::CommandDao> &cmdDao,
        const std::string &mqttIp, const std::uint32_t port,
        const std::string &deviceID);

   private:
    std::shared_ptr<db::CommandDao> cmdDao = nullptr;

    std::priority_queue<model::Command> cmdQueue;
    std::set<std::string>               queuedCmds;

    std::string         deviceID;
    network::MqttClient mqttClient;

    model::Command decodeCmdMsg(mqtt::const_message_ptr msg);

    void onCmdArrived(mqtt::const_message_ptr);
    void enqueueCommand(const model::Command &cmd);
};
};  // namespace OcelotMDM::component::service
