#pragma once

#include <curl/curl.h>
#include <curl/easy.h>

#include <cstddef>
#include <list>
#include <optional>
#include <string>

#include "http_response.hpp"

namespace OcelotMDM::component::network {
class HttpClient {
   public:
    explicit HttpClient(const std::string &baseUrl);
    ~HttpClient();

    httpResponse get(
        const std::string &path, const std::list<std::string> &header,
        const std::optional<std::string> &queryParams = std::nullopt);
    httpResponse post(
        const std::string &path, const std::list<std::string> &header,
        const std::string                &body,
        const std::optional<std::string> &queryParams = std::nullopt);

   private:
    std::string baseUrl;
    CURL       *curlHandle;

    void resetOpts();

    static curl_slist *generateHeader(const std::list<std::string> &header);
    static size_t      write_callback(
             char *ptr, size_t size, size_t nmemb, void *userdata);
    static std::string buildUrl(
        const std::string &basePath, const std::string &path,
        const std::optional<std::string> &queryParams);
};
};  // namespace OcelotMDM::component::network
