#include "command_dao.hpp"

#include <sqlite3.h>

#include <cstddef>
#include <cstdint>
#include <iostream>
#include <list>
#include <memory>
#include <optional>
#include <string>

#include "command_model.hpp"
#include "status_codes.hpp"

namespace OcelotMDM::component::db {
CommandDao::CommandDao(const std::shared_ptr<sqlite3> &dbConn) {
    this->dbConn = dbConn;
}

std::optional<bool> CommandDao::enqueueCommand(const model::Command &cmd) {
    if (this->dbConn == nullptr) {
        this->error = "db conn is not initialized";
        return std::nullopt;
    }

    std::string query(
        "INSERT INTO COMMANDS (value, action, payload, priority, queued) "
        "values (?, ?, ?, ?, ?);");

    sqlite3_stmt *stmt;
    int           ret;

    ret = sqlite3_prepare_v2(
        this->dbConn.get(), query.c_str(), -1, &stmt, nullptr);

    if (ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return std::nullopt;
    }

    sqlite3_bind_text(stmt, 1, cmd.id.c_str(), cmd.id.size(), SQLITE_STATIC);
    sqlite3_bind_text(
        stmt, 2, cmd.commandAction.c_str(), cmd.commandAction.size(),
        SQLITE_STATIC);
    sqlite3_bind_text(
        stmt, 3, cmd.payload.c_str(), cmd.payload.size(), SQLITE_STATIC);
    sqlite3_bind_int(stmt, 4, cmd.priority);
    sqlite3_bind_int(stmt, 5, true);

    sqlite3_step(stmt);

    ret = sqlite3_finalize(stmt);

    if (ret != SQLITE_DONE && ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return false;
    }

    return true;
}

std::optional<std::list<model::Command>> CommandDao::getQueuedCommands() {
    if (this->dbConn == nullptr) {
        this->error = "db conn is not initialized";
        return std::nullopt;
    }

    std::string query("select * from COMMANDS where queued=1;");

    sqlite3_stmt *stmt;
    int           ret;

    ret = sqlite3_prepare_v2(
        this->dbConn.get(), query.c_str(), -1, &stmt, nullptr);

    if (ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return std::nullopt;
    }

    std::list<model::Command> cmdList;
    while ((ret = sqlite3_step(stmt)) == SQLITE_ROW) {
        auto id = reinterpret_cast<const char *>(sqlite3_column_text(stmt, 1));
        auto action =
            reinterpret_cast<const char *>(sqlite3_column_text(stmt, 2));
        auto payload =
            reinterpret_cast<const char *>(sqlite3_column_text(stmt, 3));
        auto priority = static_cast<std::uint32_t>(sqlite3_column_int(stmt, 4));

        cmdList.emplace_back(id, action, payload, priority);
    }

    sqlite3_finalize(stmt);
    if (ret != SQLITE_DONE) {
        this->error = sqlite3_errstr(ret);
        return {};
    }

    return cmdList;
}

std::optional<bool> CommandDao::dequeCommand(const std::string &id) {
    if (this->dbConn == nullptr) {
        this->error = "db conn is not initialized";
        return std::nullopt;
    }

    std::string query("update commands set queued = 0 where value = ?");

    sqlite3_stmt *stmt;
    int           ret;

    ret = sqlite3_prepare_v2(
        this->dbConn.get(), query.c_str(), -1, &stmt, nullptr);

    if (ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return std::nullopt;
    }

    sqlite3_bind_text(stmt, 1, id.c_str(), id.size(), SQLITE_STATIC);

    sqlite3_step(stmt);

    ret = sqlite3_finalize(stmt);

    if (ret != SQLITE_DONE && ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return false;
    }

    return true;
}

std::string CommandDao::getError() {
    return this->error;
}
}  // namespace OcelotMDM::component::db
