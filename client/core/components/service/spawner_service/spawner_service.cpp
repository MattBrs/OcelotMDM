#include "spawner_service.hpp"

#include <iostream>

#include "binary_model.hpp"

namespace OcelotMDM::component::service {
SpawnerService::SpawnerService(const std::shared_ptr<db::BinaryDao> &binDao) {
    this->binDao = binDao;
}

void SpawnerService::startBinaries() {
    auto binaries = this->binDao->listBinaries();
    if (binaries.has_value()) {
        for (const auto &bin : binaries.value()) {
            std::cout << "found binary: " << bin.getPath() << std::endl;
            this->runBinary(bin.getPath());
        }
    }
}
}  // namespace OcelotMDM::component::service
