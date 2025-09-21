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
#include "commands_impl.hpp"
#include "http_client.hpp"
#include "mqtt_client.hpp"

namespace OcelotMDM::component::service {
class CommandService {
   public:
    CommandService(
        const std::shared_ptr<db::CommandDao> &cmdDao,
        const std::string &mqttIp, const std::uint32_t port,
        const std::string &httpBaseUrl, const std::string &deviceID);

    ~CommandService();

   private:
    const int DEQUEUE_INTR = 5000;

    std::shared_ptr<db::CommandDao> cmdDao = nullptr;

    std::string                         deviceID;
    std::priority_queue<model::Command> cmdQueue;
    network::MqttClient                 mqttClient;
    network::HttpClient                 httpClient;
    std::set<std::string>               queuedCmds;

    std::thread             queueTh;
    std::condition_variable queueCv;
    std::mutex              queueMtx;
    std::atomic<bool>       shouldStopTh;

    model::Command decodeCmdMsg(mqtt::const_message_ptr msg);

    void        onCmdArrived(mqtt::const_message_ptr);
    void        enqueueCommand(model::Command &cmd);
    void        queueWorker();
    std::string encodeCmd(const model::Command &cmd);
    std::optional<CommandImpl::ExecutionResult> executeCommand(
        const model::Command &cmd);
};
};  // namespace OcelotMDM::component::service
