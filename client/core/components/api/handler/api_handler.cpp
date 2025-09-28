#include "api_handler.hpp"

#include <fstream>
#include <iostream>
#include <nlohmann/json.hpp>
#include <nlohmann/json_fwd.hpp>
#include <optional>
#include <sstream>
#include <string>

#include "dto.hpp"
#include "logger.hpp"

namespace OcelotMDM::component::Api::Handler {
std::optional<component::Api::Dto::GetBinaryRes> getBinary(
    network::HttpClient *client, const component::Api::Dto::GetBinaryReq &req) {
    std::stringstream ss{""};
    ss << "?name=" << req.binaryName << "&token=" << req.otp;

    auto res = client->get("/binary/get", {}, ss.str());

    Logger::getInstance().put("Http call done");
    Logger::getInstance().put("- " + std::to_string(res.statusCode));
    if (res.error.has_value()) {
        Logger::getInstance().put("- " + res.error.value());
    }

    if (res.data.has_value()) {
        Logger::getInstance().put("- " + res.data.value());
    }

    if (res.statusCode != 200 || !res.data.has_value()) {
        return std::nullopt;
    }

    component::Api::Dto::GetBinaryRes getBinRes;
    try {
        nlohmann::json body = nlohmann::json::parse(res.data.value());
        getBinRes.name = body["binary_name"];
        getBinRes.version = body["binary_version"];
        getBinRes.binaryData = body["binary_data"];
    } catch (nlohmann::json::exception &e) {
        std::cout << "errored on res contruction: " << e.what() << std::endl;
        return std::nullopt;
    }

    return getBinRes;
}
}  // namespace OcelotMDM::component::Api::Handler
