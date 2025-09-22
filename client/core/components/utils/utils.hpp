#pragma once

#include <optional>
#include <string>
#include <vector>
namespace OcelotMDM::component::utils {
std::vector<std::string> splitString(
    const std::string &s, const std::string &delim);

std::optional<std::string> readFile(const std::string &filePath);
};  // namespace OcelotMDM::component::utils
