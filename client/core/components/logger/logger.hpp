#pragma once

#include <fstream>
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

   private:
    std::ofstream output;

    void log(const std::string &);
};
};  // namespace OcelotMDM::component
