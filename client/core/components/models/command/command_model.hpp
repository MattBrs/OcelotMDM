#pragma once

#include <cstdint>
#include <string>
namespace OcelotMDM::component::model {
struct Command {
    std::string   id;
    std::string   commandAction;
    std::string   payload;
    std::uint32_t priority;

    bool operator<(const Command &other) const {
        return priority < other.priority;
    }
};
};  // namespace OcelotMDM::component::model
