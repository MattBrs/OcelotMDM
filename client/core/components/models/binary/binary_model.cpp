#include "binary_model.hpp"

namespace OcelotMDM::component::model {
Binary::Binary(const std::string &name, const std::string &path) {
    this->name = name;
    this->path = path;
}

std::string Binary::getId() {
    return this->id;
}

std::string Binary::getName() {
    return this->name;
}

std::string Binary::getPath() {
    return this->path;
}

void Binary::setId(const std::string &id) {
    this->id = id;
}

void Binary::setName(const std::string &name) {
    this->name = name;
}

void Binary::setPath(const std::string &path) {
    this->path = path;
}

};  // namespace OcelotMDM::component::model
