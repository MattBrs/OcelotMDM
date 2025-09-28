#pragma once

#include <memory>
#include <string>

#include "binary_dao.hpp"
#include "spawner_service.hpp"

namespace OcelotMDM::component::service {
class LinuxSpawerService : public SpawnerService {
   public:
    LinuxSpawerService(const std::shared_ptr<db::BinaryDao> &binDao);

    int runBinary(const std::string &path);

   private:
};
}  // namespace OcelotMDM::component::service
