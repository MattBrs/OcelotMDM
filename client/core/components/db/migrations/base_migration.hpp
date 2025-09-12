#pragma once

#include <string>
namespace OcelotMDM::component::db {
class BaseMigration {
   public:
    BaseMigration() = default;
    virtual ~BaseMigration() = default;

    virtual std::string getMigration() = 0;

   protected:
    int version = 0;
};
}  // namespace OcelotMDM::component::db
