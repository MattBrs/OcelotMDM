#pragma once

#include <cstdint>
#include <string>
namespace OcelotMDM::component::model {
struct Command {
    std::string   id;
    std::string   commandAction;
    std::string   payload;
    std::uint32_t priority;

    Command(
        const std::string &_id, const std::string &_action,
        const std::string &_payload, const std::uint32_t _priority)
        : id(_id),
          commandAction(_action),
          payload(_payload),
          priority(_priority) {}

    bool operator<(const Command &other) const {
        return priority < other.priority;
    }
};
};  // namespace OcelotMDM::component::model
