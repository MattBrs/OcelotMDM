#include "logger.hpp"

#include <cstdio>
#include <ios>
#include <iostream>
#include <source_location>
#include <string>

namespace OcelotMDM::component {
Logger::Logger() : output("ocelot_output.txt", std::ios_base::app) {
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
    std::clog << data << std::endl;
}
};  // namespace OcelotMDM::component
