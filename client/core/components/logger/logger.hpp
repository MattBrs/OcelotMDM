#pragma once

#include <fstream>
#include <mutex>
#include <source_location>
#include <string>

namespace OcelotMDM::component {
class Logger {
   public:
    Logger();
    ~Logger();

    static Logger &getInstance();

    void put(
        const std::string &,
        const std::source_location & = std::source_location::current());
    void putError(
        const std::string &,
        const std::source_location & = std::source_location::current());

    void switchFile();

   private:
    const int MAX_LOG_SIZE = 1048576;  // 1MB

    std::mutex    logMtx;
    std::ofstream output;
    std::string   currentFileName;
    int           currentFileSize = 0;

    void        log(const std::string &);
    std::string generateFileName() const;
};
};  // namespace OcelotMDM::component
