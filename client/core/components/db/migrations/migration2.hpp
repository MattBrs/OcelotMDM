#pragma once

#include <string>

#include "base_migration.hpp"

namespace OcelotMDM::component::db {
class Migration2 : public BaseMigration {
   public:
    Migration2() {
        this->version = 2;
    }

    std::string getMigration() {
        return R"(
    CREATE TABLE IF NOT EXISTS BINARIES(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE, 
        path TEXT NOT NULL UNIQUE,
        installet_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    PRAGMA user_version = 2;)";
    }
};
}  // namespace OcelotMDM::component::db
