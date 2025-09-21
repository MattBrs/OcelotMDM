#pragma once

#include <string>
#include <vector>
namespace OcelotMDM::component::utils {
std::vector<std::string> splitString(
    const std::string &s, const std::string &delim);
};  // namespace OcelotMDM::component::utils
