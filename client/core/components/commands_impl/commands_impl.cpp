#include "commands_impl.hpp"

#include <filesystem>
#include <fstream>
#include <ios>
#include <iostream>
#include <string>
#include <vector>

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
    network::MqttClient *client, const std::string &deviceID) {
    executionResult res;

    Logger::getInstance().switchFile();

    std::string logPath{"logs"};
    std::string currentLogName = Logger::getInstance().getCurrentLogName();
    std::string topic = std::string{deviceID}.append("/logs");

    for (const auto &entry : std::filesystem::directory_iterator(logPath)) {
        auto file = entry.path();
        if (file.compare(currentLogName) == 0) {
            // skip the current file
            continue;
        }

        auto fileData = CommandImpl::readFile(file);
        client->publish(fileData, topic, 1);
        std::filesystem::remove(file);

        Logger::getInstance().put("Sent log: " + file.filename().string());
    }

    res.successful = true;
    return res;
}

std::string CommandImpl::readFile(const std::string &filePath) {
    std::ifstream ifs(
        filePath.c_str(), std::ios::in | std::ios::binary | std::ios::ate);
    auto fileSize = ifs.tellg();
    ifs.seekg(0, std::ios::beg);

    std::vector<char> bytes(fileSize);
    ifs.read(bytes.data(), fileSize);

    return std::string{bytes.data(), static_cast<unsigned long>(fileSize)};
}
};  // namespace OcelotMDM::component
