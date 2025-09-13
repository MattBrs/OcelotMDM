#pragma once

#include <cstdint>
#include <string>
namespace OcelotMDM::component::model {

class Command {
   public:
    enum class CommandStatus { Acknowledged, Completed, Errored };

    Command(
        const std::string &_id, const std::string &_action,
        const std::string &_payload, const std::uint32_t _priority,
        const bool          _onlineRequired,
        const CommandStatus _status = CommandStatus::Acknowledged,
        const std::string  &_error = "");

    bool operator<(const Command &other) const;

    std::string   getId() const;
    std::string   getAction() const;
    std::string   getPayload() const;
    std::uint32_t getPriority() const;
    std::string   getStatus() const;
    std::string   getError() const;
    bool          isOnlineRequired() const;

    void setStatus(const CommandStatus status);
    void setError(const std::string &error);

   private:
    std::string   id;
    std::string   commandAction;
    std::string   payload;
    std::uint32_t priority;
    CommandStatus status;
    std::string   errorMsg;
    bool          requireOnline;
};

};  // namespace OcelotMDM::component::model
