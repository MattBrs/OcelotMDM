#pragma once

#include <sqlite3.h>

#include <array>
#include <list>
#include <memory>
#include <string>
#include <vector>

#include "base_migration.hpp"
#include "command_dao.hpp"
#include "migration1.hpp"

namespace OcelotMDM::component::db {
class DbClient {
   public:
    explicit DbClient(const std::string &dbName);

    DbClient(const DbClient &other) = delete;
    DbClient(DbClient &&other) = delete;
    DbClient &operator=(const DbClient &other) = delete;
    DbClient &operator=(DbClient &&other) noexcept;

    std::shared_ptr<CommandDao> getCommandDao();

   private:
    std::shared_ptr<sqlite3>                    db;
    std::vector<std::unique_ptr<BaseMigration>> migrations;

    std::shared_ptr<CommandDao> commandDao = nullptr;

    int  getSchemaVersion(std::string &error);
    void runMigrations(const int currentVersion);
    bool executeString(const std::string &query, std::string &error);
};
};  // namespace OcelotMDM::component::db
