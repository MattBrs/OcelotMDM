#include "commands_impl.hpp"

#include <iostream>

#include "logger.hpp"

namespace OcelotMDM::component {
CommandImpl::ExecutionResult CommandImpl::installBinary(
    network::HttpClient *client, const std::string &url) {
    executionResult res;

    // fetch binary with httpClient
    // install to binaries directory
    // return the binary info so that the application can be managed

    Logger::getInstance().put(
        "shound run command install_command for url:" + url);

    res.props.error = "not implemented yet";
    return res;
}

CommandImpl::ExecutionResult CommandImpl::sendLogs(
    network::MqttClient *client) {
    executionResult res;

    // read the log file
    // send the logs through mqtt <device_id>/logs topic
    // delete the log file since it's synced

    Logger::getInstance().put("shound run command send_logs");

    res.props.error = "not implemented yet";
    return res;
}
};  // namespace OcelotMDM::component
