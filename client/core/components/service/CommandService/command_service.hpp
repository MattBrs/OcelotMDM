#pragma once

#include <mqtt/message.h>

#include <atomic>
#include <condition_variable>
#include <cstdint>
#include <memory>
#include <mutex>
#include <queue>
#include <set>
#include <string>
#include <thread>

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
    ~CommandService();

   private:
    const int DEQUEUE_INTR = 5000;

    std::shared_ptr<db::CommandDao> cmdDao = nullptr;

    std::priority_queue<model::Command> cmdQueue;
    std::set<std::string>               queuedCmds;

    std::string         deviceID;
    network::MqttClient mqttClient;

    std::thread             queueTh;
    std::condition_variable queueCv;
    std::mutex              queueMtx;
    std::atomic<bool>       shouldStopTh;

    model::Command decodeCmdMsg(mqtt::const_message_ptr msg);

    void onCmdArrived(mqtt::const_message_ptr);
    void enqueueCommand(const model::Command &cmd);
    void queueWorker();
};
};  // namespace OcelotMDM::component::service
