#pragma once

#include <optional>
#include <string>
namespace OcelotMDM::component::network {
struct httpResponse {
    std::optional<std::string> data = std::nullopt;
    std::optional<std::string> error = std::nullopt;
    int                        statusCode = -1;
};
};  // namespace OcelotMDM::component::network
