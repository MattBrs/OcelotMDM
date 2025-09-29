#pragma once

#include <memory>
#include <string>

#include "binary_dao.hpp"
#include "http_client.hpp"
#include "log_streamer.hpp"
#include "mqtt_client.hpp"
#include "timer.hpp"

namespace OcelotMDM::component {
class CommandImpl {
   public:
    typedef struct cmdResProps {
        std::string applicationName;
        std::string applicationPath;
        std::string error;
        std::string data;
    } CmdResProps;

    typedef struct executionResult {
        CmdResProps props;
        bool        successful = false;
    } ExecutionResult;

    /**
     *  If successful, returns the applicationName and applicationPath in the
     * props, otherwise the error
     */
    static ExecutionResult installBinary(
        const std::shared_ptr<db::BinaryDao> &binDao,
        std::shared_ptr<network::HttpClient> &httpClient,
        const std::string &name, const std::string &otp);

    /**
     *  If successful, starts the target binary on the device, otherwise the
     * cause of error
     */
    static ExecutionResult startBinary(
        const std::shared_ptr<db::BinaryDao> &binDao, const std::string &name);

    /**
     *  If successful, starts the target binary on the device, otherwise the
     * cause of error
     */
    static ExecutionResult uninstallBinary(
        const std::shared_ptr<db::BinaryDao> &binDao, const std::string &name);

    /**
     *  If successful, starts the target binary on the device, otherwise the
     * cause of error
     */
    static ExecutionResult listBinaries(
        const std::shared_ptr<db::BinaryDao> &binDao);

    /**
     *  If successful, returns logData inside props, otherwise the error
     */
    static ExecutionResult sendLogs(
        const std::shared_ptr<network::MqttClient> &client,
        const std::string                          &deviceID);

    /**
     *  If successful, enables the streaming of all the logs to the topic
     * "live-logs", otherwise error
     */
    static ExecutionResult enableLiveLogging(
        const std::shared_ptr<LogStreamer> &logStreamer,
        const std::shared_ptr<Timer>       &timer);
    /**
     *  If successful, disables the streaming of logs
     */
    static ExecutionResult disableLiveLogging(
        const std::shared_ptr<LogStreamer> &logStreamer,
        const std::shared_ptr<Timer>       &timer);

   private:
};
}  // namespace OcelotMDM::component
