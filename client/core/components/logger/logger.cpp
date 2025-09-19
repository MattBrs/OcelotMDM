#include "logger.hpp"

#include <chrono>
#include <cstdio>
#include <cstdlib>
#include <ctime>
#include <filesystem>
#include <fstream>
#include <ios>
#include <iostream>
#include <source_location>
#include <string>

namespace OcelotMDM::component {
Logger::Logger() {
    std::filesystem::create_directory("logs");

    this->currentFileName = this->generateFileName();
    this->output = std::ofstream(currentFileName, std::ios_base::app);

    std::clog.rdbuf(this->output.rdbuf());
}

Logger::~Logger() {
    std::clog.flush();
    output.flush();
    std::clog.rdbuf(nullptr);
    this->output.close();
}

Logger &Logger::getInstance() {
    static Logger logger;
    return logger;
}

void Logger::put(
    const std::string &data, const std::source_location &location) {
    std::string completeLog;

    completeLog.append(location.file_name())
        .append(":")
        .append(std::to_string(location.line()))
        .append("    ")
        .append(data);

    this->log(completeLog);
}

void Logger::putError(
    const std::string &data, const std::source_location &location) {
    std::string completeLog;

    completeLog.append(location.file_name())
        .append(":")
        .append(std::to_string(location.line()))
        .append("    ")
        .append("Error: ")
        .append(data);

    this->log(completeLog);
}

void Logger::log(const std::string &data) {
    this->logMtx.lock();

    this->currentFileSize += data.size();
    std::clog << data << std::endl;

    this->logMtx.unlock();

    if (this->currentFileSize > this->MAX_LOG_SIZE) {
        this->switchFile();
    }
}

void Logger::switchFile() {
    this->logMtx.lock();
    std::clog.rdbuf(nullptr);

    this->output.close();
    this->currentFileName = this->generateFileName();
    this->output = std::ofstream(currentFileName, std::ios_base::app);

    std::clog.rdbuf(this->output.rdbuf());

    this->logMtx.unlock();
}

std::string Logger::generateFileName() const {
    auto        now = std::chrono::system_clock::now();
    std::time_t epochTime = std::chrono::system_clock::to_time_t(now);

    auto fileName = std::string{"logs/ocelot_logs_"}.append(
        std::to_string(epochTime).append(".txt"));

    return fileName;
}

};  // namespace OcelotMDM::component
