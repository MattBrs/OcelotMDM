#include "uptime_service.hpp"

#include <ifaddrs.h>
#include <netdb.h>
#include <sys/socket.h>

#include <chrono>
#include <ctime>
#include <iostream>
#include <memory>
#include <mutex>
#include <optional>
#include <sstream>
#include <string>

#include "mqtt_client.hpp"

namespace OcelotMDM::component::service {
UptimeService::UptimeService(
    const std::shared_ptr<network::MqttClient> &client,
    const std::string                          &deviceID)
    : deviceID(deviceID), mqttClient(client) {
    this->shouldStop.store(false);
    this->workerTh = std::thread(&UptimeService::workerFunction, this);
}

UptimeService::~UptimeService() {
    this->shouldStop.store(true);
    this->workerCv.notify_one();

    if (this->workerTh.joinable()) {
        this->workerTh.join();
    }
}

void UptimeService::workerFunction() {
    while (!this->shouldStop.load()) {
        auto vpnIp = this->getVpnAddress();

        if (vpnIp.has_value()) {
            std::stringstream ss;
            auto              now = time(nullptr);
            auto              localNow = std::mktime(std::localtime(&now));

            ss << localNow << " " << vpnIp.value();
            this->mqttClient->publish(
                ss.str(), this->deviceID + "/online", 1, true);
        }

        std::unique_lock<std::mutex> lock(this->workerMtx);
        this->workerCv.wait_for(
            lock, std::chrono::milliseconds(this->UPDATE_FREQUENCY),
            [this]() { return this->shouldStop.load(); });
    }
}

std::optional<std::string> UptimeService::getVpnAddress() {
    struct ifaddrs            *addresses;
    std::optional<std::string> ip = std::nullopt;

    if (getifaddrs(&addresses) == -1) {
        return ip;
    }

    struct ifaddrs *address = addresses;
    while (address) {
        sa_family_t family;
        std::string iName = address->ifa_name;
        auto        interfaceInfo = address->ifa_addr;

        if (interfaceInfo == nullptr) {
            family = AF_INET;
        } else {
            family = address->ifa_addr->sa_family;
        }

        if (family == AF_INET && iName.find("tun") != std::string::npos) {
            char      ap[100];
            const int familySize = family == AF_INET
                                       ? sizeof(struct sockaddr_in)
                                       : sizeof(struct sockaddr_in6);
            getnameinfo(
                address->ifa_addr, familySize, ap, sizeof(ap), 0, 0,
                NI_NUMERICHOST);

            ip = ap;
        }

        address = address->ifa_next;
    }

    freeifaddrs(addresses);
    return ip;
}
}  // namespace OcelotMDM::component::service
