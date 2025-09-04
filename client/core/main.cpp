#include <iostream>

#include "http_client.hpp"

using HttpClient = OcelotMDM::component::network::HttpClient;

int main() {
    HttpClient httpClient("https://httpbin.org/ip");
    auto       res = httpClient.get("", {});

    if (res.data.has_value()) {
        std::cout << "response body: " << res.data.value();
        std::cout << "response code: " << res.statusCode << "\n";
        std::cout << std::flush;
    } else if (res.error.has_value()) {
        std::cout << "error: " << res.error.value() << "\n";
        std::cout << std::flush;
    }

    std::cout << "Hello!" << std::endl;
}
