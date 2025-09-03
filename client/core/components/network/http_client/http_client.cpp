#include "http_client.hpp"

#include <curl/curl.h>
#include <curl/easy.h>

#include <list>
#include <optional>
#include <string>

namespace OcelotMDM::component::network {
HttpClient::HttpClient(const std::string &baseUrl) {
    this->baseUrl = baseUrl;

    this->curlHandle = curl_easy_init();
};

httpResponse HttpClient::get(
    const std::string &path, const std::list<std::string> &header) {
    if (!this->curlHandle) {
        // client not properly initialized
        return {std::nullopt, "curl handle not initialized", -1};
    }

    auto url = this->buildUrl(path, std::nullopt);
    curl_easy_setopt(this->curlHandle, CURLOPT_URL, url.c_str());
    curl_easy_setopt(this->curlHandle, CURLOPT_FOLLOWLOCATION, url.c_str());

    return {std::nullopt, std::nullopt, 200};
}

httpResponse HttpClient::post(
    const std::string &path, const std::list<std::string> &header,
    const std::string &body) {
    if (!this->curlHandle) {
        // client not properly initialized
        return {std::nullopt, "curl handle not initialized", -1};
    }

    auto url = this->buildUrl(path, std::nullopt);
    curl_easy_setopt(this->curlHandle, CURLOPT_URL, url.c_str());
    curl_easy_setopt(this->curlHandle, CURLOPT_FOLLOWLOCATION, url.c_str());
    // curl_easy_setopt(this->curlHandle, , opt, param) // set post body

    return {std::nullopt, std::nullopt, 200};
}

std::string HttpClient::buildUrl(
    const std::string &path, const std::optional<std::string> &queryParams) {
    auto completePath = this->baseUrl;
    if (path.starts_with("/")) {
        completePath.append(path.begin() + 1, path.end());
    } else {
        completePath.append(path);
    }

    if (!queryParams.has_value()) {
        return completePath;
    }

    if (!queryParams->starts_with("?")) {
        completePath.append("?");
    }

    completePath.append(queryParams.value());
    return completePath;
}
};  // namespace OcelotMDM::component::network
