#pragma once

#include <string>

#include "http_client.hpp"
#include "mqtt_client.hpp"

namespace OcelotMDM::component {
class CommandImpl {
   public:
    typedef struct cmdResProps {
        std::string applicationName;
        std::string applicationPath;
        std::string logData;
        std::string error;
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
        network::HttpClient *client, const std::string &url);

    /**
     *  If successful, returns logData inside props, otherwise the error
     */
    static ExecutionResult sendLogs(network::MqttClient *client);
};
}  // namespace OcelotMDM::component
