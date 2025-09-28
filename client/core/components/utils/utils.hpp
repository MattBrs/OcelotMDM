#pragma once

#include <optional>
#include <string>
#include <vector>

namespace OcelotMDM::component::utils {
typedef unsigned char BYTE;

std::string base64_encode(BYTE const *buf, unsigned int bufLen);

std::vector<BYTE> base64_decode(std::string const &);

std::vector<std::string> splitString(
    const std::string &s, const std::string &delim);

std::optional<std::string> readFile(const std::string &filePath);
};  // namespace OcelotMDM::component::utils
