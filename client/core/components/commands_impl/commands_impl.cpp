#include "commands_impl.hpp"

#include <filesystem>
#include <fstream>
#include <ios>
#include <iostream>
#include <memory>
#include <string>
#include <vector>

#include "api_handler.hpp"
#include "binary_dao.hpp"
#include "dto.hpp"
#include "logger.hpp"
#include "mqtt_client.hpp"
#include "timer.hpp"
#include "utils.hpp"

namespace OcelotMDM::component {
CommandImpl::ExecutionResult CommandImpl::installBinary(
    const std::shared_ptr<db::BinaryDao> &binDao,
    std::shared_ptr<network::HttpClient> &httpClient, const std::string &name,
    const std::string &otp) {
    executionResult res;
    std::filesystem::create_directory("bin");
    if (httpClient == nullptr) {
        res.props.error = "http client not initialized";
        return res;
    }

    Logger::getInstance().put("inside install binary");
    Logger::getInstance().put(name);
    Logger::getInstance().put(otp);
    auto httpRes = Api::Handler::getBinary(
        httpClient.get(),
        component::Api::Dto::GetBinaryReq{.binaryName = name, .otp = otp});
    Logger::getInstance().put("after http call");

    if (!httpRes.has_value()) {
        Logger::getInstance().put("error on binary fetch");
        res.props.error = "error on binary fetch";
        return res;
    }

    auto binary = utils::base64_decode(httpRes.value().binaryData);

    auto          appPath = std::string{"bin/"}.append(name);
    std::ofstream out(appPath, std::ios::out | std::ios::binary);

    out.write(reinterpret_cast<const char *>(binary.data()), binary.size());
    out.close();

    std::filesystem::permissions(
        appPath,
        std::filesystem::perms::owner_all | std::filesystem::perms::group_all,
        std::filesystem::perm_options::add);

    auto insertRes = binDao->addBinary(name, appPath);
    if (!insertRes.has_value() || !insertRes.value()) {
        res.successful = false;
        res.props.error = binDao->getError();
        return res;
    }

    res.successful = true;
    res.props.applicationName = name;
    res.props.applicationPath = appPath;

    return res;
}

CommandImpl::ExecutionResult CommandImpl::sendLogs(
    const std::shared_ptr<network::MqttClient> &client,
    const std::string                          &deviceID) {
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

        auto fileData = utils::readFile(file);
        if (!fileData.has_value()) {
            continue;
        }

        client->publish(fileData.value(), topic, 1);
        std::filesystem::remove(file);

        Logger::getInstance().put("Sent log: " + file.filename().string());
    }

    res.successful = true;
    return res;
}

CommandImpl::ExecutionResult CommandImpl::enableLiveLogging(
    const std::shared_ptr<LogStreamer> &logStreamer,
    const std::shared_ptr<Timer>       &timer) {
    ExecutionResult res;
    if (logStreamer->isRunning()) {
        res.successful = false;
        res.props.error = "live logging already enabled";
        return res;
    }

    auto streamerQueue = logStreamer->getQueue();
    Logger::getInstance().registerQueue(streamerQueue);
    logStreamer->run();

    timer->start(
        [logStreamer]() {
            Logger::getInstance().registerQueue(nullptr);

            if (logStreamer != nullptr) {
                logStreamer->stop();
            }
        },
        10 * 60 * 1000, false);

    res.successful = true;
    return res;
}
};  // namespace OcelotMDM::component
