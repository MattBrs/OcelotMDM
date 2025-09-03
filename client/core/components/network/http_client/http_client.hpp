#pragma once

#include <curl/curl.h>
#include <curl/easy.h>

#include <list>
#include <optional>
#include <string>

#include "http_response.hpp"

namespace OcelotMDM::component::network {
class HttpClient {
   public:
    explicit HttpClient(const std::string &baseUrl);

    httpResponse get(
        const std::string &path, const std::list<std::string> &header);
    httpResponse post(
        const std::string &path, const std::list<std::string> &header,
        const std::string &body);

   private:
    std::string baseUrl;
    CURL       *curlHandle;

    std::string buildUrl(
        const std::string &path, const std::optional<std::string> &queryParams);
};
};  // namespace OcelotMDM::component::network
