#pragma once

#include <memory>
#include <string>

#include "binary_dao.hpp"
namespace OcelotMDM::component::service {
class SpawnerService {
   public:
    SpawnerService(const std::shared_ptr<db::BinaryDao> &binDao);
    virtual ~SpawnerService() = default;

    virtual int runBinary(const std::string &path) = 0;

   protected:
    std::shared_ptr<db::BinaryDao> binDao = nullptr;

    void startBinaries();
};
}  // namespace OcelotMDM::component::service
