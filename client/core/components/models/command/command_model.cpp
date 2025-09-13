#include "command_model.hpp"

namespace OcelotMDM::component::model {
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

void Command::setStatus(const CommandStatus status) {
    this->status = status;
}

void Command::setError(const std::string &error) {
    this->errorMsg = error;
}

bool Command::operator<(const Command &other) const {
    return priority < other.priority;
}
};  // namespace OcelotMDM::component::model
