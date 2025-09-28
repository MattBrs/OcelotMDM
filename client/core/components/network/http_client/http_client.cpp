#include "http_client.hpp"

#include <curl/curl.h>
#include <curl/easy.h>

#include <fstream>
#include <iostream>
#include <list>
#include <memory>
#include <optional>
#include <sstream>
#include <string>

#include "http_response.hpp"
#include "logger.hpp"

namespace OcelotMDM::component::network {
HttpClient::HttpClient(const std::string &baseUrl) {
    this->baseUrl = baseUrl;

    this->curlHandle = curl_easy_init();

    curl_easy_setopt(
        this->curlHandle, CURLOPT_WRITEFUNCTION, &HttpClient::write_callback);
};

HttpClient::~HttpClient() {
    curl_easy_cleanup(this->curlHandle);
}

httpResponse HttpClient::get(
    const std::string &path, const std::list<std::string> &header,
    const std::optional<std::string> &queryParams) {
    if (!this->curlHandle) {
        // client not properly initialized
        return {std::nullopt, "curl handle not initialized", -1};
    }

    std::string resData;
    char        errorBuffer[CURL_ERROR_SIZE];
    auto        curlHeader = HttpClient::generateHeader(header);
    auto        url = this->buildUrl(this->baseUrl, path, queryParams);

    Logger::getInstance().put("full http url: " + url);

    errorBuffer[0] = 0;

    curl_easy_setopt(this->curlHandle, CURLOPT_HTTPGET, 1L);
    curl_easy_setopt(this->curlHandle, CURLOPT_TIMEOUT, 15L);
    curl_easy_setopt(this->curlHandle, CURLOPT_URL, url.c_str());
    curl_easy_setopt(this->curlHandle, CURLOPT_HTTPHEADER, curlHeader);
    curl_easy_setopt(this->curlHandle, CURLOPT_ERRORBUFFER, errorBuffer);
    curl_easy_setopt(this->curlHandle, CURLOPT_FOLLOWLOCATION, url.c_str());
    curl_easy_setopt(this->curlHandle, CURLOPT_WRITEDATA, &resData);
    curl_easy_setopt(this->curlHandle, CURLOPT_SSL_VERIFYPEER, 0);
    curl_easy_setopt(this->curlHandle, CURLOPT_SSL_VERIFYHOST, 0);

    auto res = curl_easy_perform(this->curlHandle);

    httpResponse response;

    switch (res) {
        case CURLcode::CURLE_OK:
            response.data = resData;
            curl_easy_getinfo(
                this->curlHandle, CURLINFO_RESPONSE_CODE, &response.statusCode);
            break;
        default:
            response.error = errorBuffer;
    }

    curl_slist_free_all(curlHeader);

    return response;
}

httpResponse HttpClient::post(
    const std::string &path, const std::list<std::string> &header,
    const std::string &body, const std::optional<std::string> &queryParams) {
    if (!this->curlHandle) {
        // client not properly initialized
        return {std::nullopt, "curl handle not initialized", -1};
    }

    std::string resData;
    char        errorBuffer[CURL_ERROR_SIZE];
    auto        curlHeader = HttpClient::generateHeader(header);
    auto        url = this->buildUrl(this->baseUrl, path, queryParams);

    errorBuffer[0] = 0;

    curl_easy_setopt(this->curlHandle, CURLOPT_POST, 1L);
    curl_easy_setopt(this->curlHandle, CURLOPT_TIMEOUT, 15L);
    curl_easy_setopt(this->curlHandle, CURLOPT_URL, url.c_str());
    curl_easy_setopt(this->curlHandle, CURLOPT_WRITEDATA, &resData);
    curl_easy_setopt(this->curlHandle, CURLOPT_HTTPHEADER, curlHeader);
    curl_easy_setopt(this->curlHandle, CURLOPT_ERRORBUFFER, errorBuffer);
    curl_easy_setopt(this->curlHandle, CURLOPT_POSTFIELDS, body.c_str());
    curl_easy_setopt(this->curlHandle, CURLOPT_POSTFIELDSIZE, body.size());
    curl_easy_setopt(this->curlHandle, CURLOPT_FOLLOWLOCATION, url.c_str());

    auto res = curl_easy_perform(this->curlHandle);

    httpResponse response;
    switch (res) {
        case CURLcode::CURLE_OK:
            response.data = resData;
            curl_easy_getinfo(
                this->curlHandle, CURLINFO_RESPONSE_CODE, &response.statusCode);
        default:
            response.error = errorBuffer;
    }

    curl_slist_free_all(curlHeader);
    return response;
}

std::string HttpClient::buildUrl(
    const std::string &basePath, const std::string &path,
    const std::optional<std::string> &queryParams) {
    auto completePath = basePath;
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

void HttpClient::resetOpts() {
    curl_easy_setopt(this->curlHandle, CURLOPT_POST, 0L);
    curl_easy_setopt(this->curlHandle, CURLOPT_POSTFIELDS, "");
    curl_easy_setopt(this->curlHandle, CURLOPT_POSTFIELDSIZE, 0);
    curl_easy_setopt(this->curlHandle, CURLOPT_HTTPGET, 0L);

    struct curl_slist *voidHeader = nullptr;
    voidHeader = curl_slist_append(voidHeader, "");
    curl_easy_setopt(this->curlHandle, CURLOPT_HTTPHEADER, voidHeader);

    curl_slist_free_all(voidHeader);
}

curl_slist *HttpClient::generateHeader(const std::list<std::string> &header) {
    struct curl_slist *curlHeader = nullptr;
    for (auto &item : header) {
        curlHeader = curl_slist_append(curlHeader, item.c_str());
    }

    return curlHeader;
}

size_t HttpClient::write_callback(
    char *ptr, size_t size, size_t nmemb, void *userdata) {
    auto         realSize = size * nmemb;
    std::string &data = *static_cast<std::string *>(userdata);
    data.append(static_cast<char *>(ptr), realSize);
    return realSize;
}
};  // namespace OcelotMDM::component::network
