#pragma once

#include <sqlite3.h>

#include <memory>
#include <optional>
#include <string>

#include "command_model.hpp"
#include "status_codes.hpp"

namespace OcelotMDM::component::db {
class CommandDao {
   public:
    explicit CommandDao(const std::shared_ptr<sqlite3> &dbConn);

    std::optional<bool> enqueueCommand(const model::Command &cmd);

    std::string getError();

   private:
    std::string              error;
    std::shared_ptr<sqlite3> dbConn;
};
};  // namespace OcelotMDM::component::db
