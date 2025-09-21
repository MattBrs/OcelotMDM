#include "commands_impl.hpp"

#include <filesystem>
#include <fstream>
#include <ios>
#include <iostream>
#include <string>
#include <vector>

#include "api_handler.hpp"
#include "dto.hpp"
#include "logger.hpp"

namespace OcelotMDM::component {
CommandImpl::ExecutionResult CommandImpl::installBinary(
    network::HttpClient *client, const std::string &name,
    const std::string &otp) {
    executionResult res;
    std::filesystem::create_directory("bin");
    if (client == nullptr) {
        res.props.error = "http client not initialized";
        return res;
    }

    Logger::getInstance().put("inside install binary");
    Logger::getInstance().put(name);
    Logger::getInstance().put(otp);
    auto httpRes = Api::Handler::getBinary(
        client,
        component::Api::Dto::GetBinaryReq{.binaryName = name, .otp = otp});
    Logger::getInstance().put("after http call");

    if (!httpRes.has_value()) {
        Logger::getInstance().put("error on binary fetch");
        res.props.error = "error on binary fetch";
        return res;
    }

    auto          appPath = std::string{"bin/"}.append(name);
    std::ofstream out(appPath);
    


    out << httpRes.value().binaryData;
    out.close();

    res.successful = true;
    res.props.applicationName = name;
    res.props.applicationPath = appPath;

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
