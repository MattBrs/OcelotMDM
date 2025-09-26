#pragma once

#include <mqtt/message.h>

#include <atomic>
#include <condition_variable>
#include <memory>
#include <mutex>
#include <queue>
#include <set>
#include <string>
#include <thread>

#include "binary_dao.hpp"
#include "command_dao.hpp"
#include "command_model.hpp"
#include "commands_impl.hpp"
#include "http_client.hpp"
#include "log_streamer.hpp"
#include "mqtt_client.hpp"
#include "timer.hpp"

namespace OcelotMDM::component::service {
class CommandService {
   public:
    CommandService(
        const std::shared_ptr<db::CommandDao>      &cmdDao,
        const std::shared_ptr<db::BinaryDao>       &binDao,
        const std::shared_ptr<network::MqttClient> &mqttClient,
        const std::string &httpBaseUrl, const std::string &deviceID);

    ~CommandService();

   private:
    const int DEQUEUE_INTR = 5000;

    std::shared_ptr<Timer>               timer = nullptr;
    std::shared_ptr<db::CommandDao>      cmdDao = nullptr;
    std::shared_ptr<db::BinaryDao>       binDao = nullptr;
    std::shared_ptr<network::MqttClient> mqttClient = nullptr;
    std::shared_ptr<LogStreamer>         logStreamer = nullptr;

    std::string                         deviceID;
    std::priority_queue<model::Command> cmdQueue;
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
