#pragma once

#include <string>

#include "base_migration.hpp"

namespace OcelotMDM::component::db {
class Migration1 : public BaseMigration {
   public:
    Migration1() {
        this->version = 1;
    }

    std::string getMigration() {
        return R"(
    CREATE TABLE IF NOT EXISTS COMMANDS(
        id INTEGER PRIMARY KEY,
        value CHAR(200) NOT NULL UNIQUE,
        action CHAR(200) NOT NULL,
        payload char(500) NOT NULL,
        priority INTEGER NOT NULL,
        queued INTEGER NOT NULL default 1,
        require_online INTEGER NOT NULL default 0
    );

    PRAGMA user_version = 1;)";
    }
};
}  // namespace OcelotMDM::component::db
