#include "utils.hpp"

#include <cstddef>
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
};  // namespace OcelotMDM::component::utils
