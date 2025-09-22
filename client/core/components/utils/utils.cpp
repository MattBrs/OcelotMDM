#include "utils.hpp"

#include <cstddef>
#include <fstream>
#include <optional>
#include <string>
#include <vector>

namespace OcelotMDM::component::utils {
std::vector<std::string> splitString(
    const std::string &s, const std::string &delim) {
    std::size_t              pos;
    std::string              og{s};
    std::string              token;
    std::vector<std::string> tokens;

    while ((pos = og.find(delim)) != std::string::npos) {
        tokens.emplace_back(og.substr(0, pos));
        og.erase(0, pos + delim.length());
    }
    tokens.emplace_back(og);

    return tokens;
}

std::optional<std::string> readFile(const std::string &filePath) {
    std::ifstream ifs(
        filePath.c_str(), std::ios::in | std::ios::binary | std::ios::ate);

    if (!ifs.is_open()) {
        return std::nullopt;
    }

    auto fileSize = ifs.tellg();
    ifs.seekg(0, std::ios::beg);

    std::vector<char> bytes(fileSize);
    ifs.read(bytes.data(), fileSize);

    return std::string{bytes.data(), static_cast<unsigned long>(fileSize)};
}
};  // namespace OcelotMDM::component::utils
