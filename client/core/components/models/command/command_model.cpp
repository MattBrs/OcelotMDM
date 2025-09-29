#include "command_model.hpp"

namespace OcelotMDM::component::model {
Command::Command(
    const std::string &_id, const std::string &_action,
    const std::string &_payload, const std::uint32_t _priority,
    const bool _onlineRequired, const CommandStatus _status,
    const std::string &_error)
    : id(_id),
      commandAction(_action),
      payload(_payload),
      priority(_priority),
      status(_status),
      errorMsg(_error),
      requireOnline(_onlineRequired) {}

std::string Command::getId() const {
    return this->id;
}

std::string Command::getAction() const {
    return this->commandAction;
}

std::string Command::getPayload() const {
    return this->payload;
}

std::uint32_t Command::getPriority() const {
    return this->priority;
}

std::string Command::getStatus() const {
    switch (this->status) {
        case CommandStatus::Acknowledged:
            return "acknowledged";
        case CommandStatus::Completed:
            return "completed";
        case CommandStatus::Errored:
            return "errored";
        default:
            return "acknowledged";
    }
}

std::string Command::getError() const {
    return this->errorMsg;
}

bool Command::isOnlineRequired() const {
    return this->requireOnline;
}

std::string Command::getData() const {
    return this->data;
}

void Command::setStatus(const CommandStatus status) {
    this->status = status;
}

void Command::setError(const std::string &error) {
    this->errorMsg = error;
}

void Command::setData(const std::string &data) {
    this->data = data;
}

bool Command::operator<(const Command &other) const {
    return priority < other.priority;
}
};  // namespace OcelotMDM::component::model
