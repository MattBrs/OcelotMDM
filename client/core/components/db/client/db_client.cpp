#include "db_client.hpp"

#include <sqlite3.h>

#include <cstdlib>
#include <iostream>
#include <memory>
#include <string>

#include "binary_dao.hpp"
#include "command_dao.hpp"
#include "migration1.hpp"
#include "migration2.hpp"

namespace OcelotMDM::component::db {
DbClient::DbClient(const std::string &dbName) {
    sqlite3 *dbInstance;
    int      result = sqlite3_open_v2(
        dbName.c_str(), &dbInstance,
        SQLITE_OPEN_READWRITE | SQLITE_OPEN_FULLMUTEX | SQLITE_OPEN_CREATE,
        nullptr);

    if (result) {
        std::cout << "erro while opening db: " << sqlite3_errstr(result)
                  << std::endl;
        exit(1);
    }

    this->db = std::shared_ptr<sqlite3>(
        dbInstance, [](sqlite3 *db) { sqlite3_close(db); });

    migrations.emplace_back(std::make_unique<Migration1>());
    migrations.emplace_back(std::make_unique<Migration2>());

    std::string err;
    int         schemaVersion = getSchemaVersion(err);
    std::cout << "schema version: " << schemaVersion << std::endl;
    if (schemaVersion < this->migrations.size()) {
        this->runMigrations(schemaVersion);
    }

    this->commandDao = std::make_shared<CommandDao>(this->db);
}

DbClient &DbClient::operator=(DbClient &&other) noexcept {
    // here we make the copies for the DAOs
    return *this;
}

int DbClient::getSchemaVersion(std::string &error) {
    if (this->db == nullptr) {
        error = "dbconn is not instantiated!";
        return false;
    }

    char       *errMsg;
    std::string select("PRAGMA user_version;");

    int schemaVersion = -1;
    int res = sqlite3_exec(
        this->db.get(), select.c_str(),
        [](void *version, int argc, char **argv, char **) {
            if (argc > 0) {
                // creates reference to the first parameter of the function
                int *t = static_cast<int *>(version);

                // saves the version on an aux variable
                int tempVer = std::stoi({argv[0]});

                // overwrites the first parameter of the function (that  it's
                // also the outside parameter given to the function).
                *t = tempVer;
            }

            return 0;
        },
        &schemaVersion, &errMsg);

    if (res) {
        error = errMsg;
        sqlite3_free(errMsg);
        return -1;
    }

    return schemaVersion;
}

void DbClient::runMigrations(const int currentVersion) {
    for (int i = currentVersion; i < this->migrations.size(); ++i) {
        std::string error;
        auto        res =
            this->executeString(this->migrations[i]->getMigration(), error);

        if (!res) {
            std::cout << "error on running migration: " << error << std::endl;
            exit(1);
        }

        std::cout << "migration " << i + 1 << " completed successfully"
                  << std::endl;
    }
}

bool DbClient::executeString(const std::string &query, std::string &error) {
    char *errMsg;
    sqlite3_exec(this->db.get(), "BEGIN;", nullptr, nullptr, nullptr);
    int ret =
        sqlite3_exec(this->db.get(), query.c_str(), nullptr, nullptr, &errMsg);

    if (ret) {
        if (errMsg != nullptr) {
            error = errMsg;
            sqlite3_free(errMsg);
        }

        sqlite3_exec(this->db.get(), "ROLLBACK;", nullptr, nullptr, nullptr);
        return false;
    }

    sqlite3_exec(this->db.get(), "COMMIT;", nullptr, nullptr, nullptr);
    return true;
}

std::shared_ptr<CommandDao> DbClient::getCommandDao() {
    return this->commandDao;
}

std::shared_ptr<BinaryDao> DbClient::getBinaryDao() {
    return this->binaryDao;
}
};  // namespace OcelotMDM::component::db
