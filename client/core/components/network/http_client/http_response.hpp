#pragma once

#include <optional>
#include <string>
namespace OcelotMDM::component::network {
struct httpResponse {
    std::optional<std::string> data;
    std::optional<std::string> error;
    int                        statusCode;
};
};  // namespace OcelotMDM::component::network
