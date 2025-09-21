#pragma once

#include <optional>

#include "dto.hpp"
#include "http_client.hpp"

namespace OcelotMDM::component::Api::Handler {
std::optional<component::Api::Dto::GetBinaryRes> getBinary(
    network::HttpClient *client, const component::Api::Dto::GetBinaryReq &req);
}  // namespace OcelotMDM::component::Api::Handler
