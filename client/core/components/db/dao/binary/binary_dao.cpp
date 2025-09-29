#include "binary_dao.hpp"

#include <sqlite3.h>

#include <memory>
#include <optional>
#include <string>
#include <vector>

#include "binary_model.hpp"

namespace OcelotMDM::component::db {
BinaryDao::BinaryDao(const std::shared_ptr<sqlite3> &dbConn) {
    this->dbConn = dbConn;
}

std::optional<bool> BinaryDao::addBinary(
    const std::string &name, const std::string &path) {
    if (this->dbConn == nullptr) {
        this->error = "db conn is not initialized";
        return std::nullopt;
    }

    std::string query("INSERT INTO BINARIES(name, path) values (?, ?);");

    sqlite3_stmt *stmt;
    int           ret;

    ret = sqlite3_prepare_v2(
        this->dbConn.get(), query.c_str(), -1, &stmt, nullptr);

    if (ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return std::nullopt;
    }

    sqlite3_bind_text(stmt, 1, name.c_str(), name.size(), SQLITE_STATIC);
    sqlite3_bind_text(stmt, 2, path.c_str(), path.size(), SQLITE_STATIC);

    sqlite3_step(stmt);
    ret = sqlite3_finalize(stmt);

    if (ret != SQLITE_DONE && ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        if (ret == SQLITE_CONSTRAINT) {
            return std::nullopt;
        }

        return false;
    }

    return true;
}

std::optional<bool> BinaryDao::removeBinary(const std::string &name) {
    if (this->dbConn == nullptr) {
        this->error = "db conn is not initialized";
        return std::nullopt;
    }

    std::string query("delete from BINARIES where name like ?;");

    sqlite3_stmt *stmt;
    int           ret;

    ret = sqlite3_prepare_v2(
        this->dbConn.get(), query.c_str(), -1, &stmt, nullptr);

    if (ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return std::nullopt;
    }

    sqlite3_bind_text(stmt, 1, name.c_str(), name.size(), SQLITE_STATIC);
    sqlite3_step(stmt);

    ret = sqlite3_finalize(stmt);

    if (ret != SQLITE_DONE && ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return false;
    }

    return true;
}

std::optional<std::vector<model::Binary>> BinaryDao::listBinaries() {
    if (this->dbConn == nullptr) {
        this->error = "db conn is not initialized";
        return std::nullopt;
    }

    std::string query("select * from BINARIES;");

    sqlite3_stmt *stmt;
    int           ret;

    ret = sqlite3_prepare_v2(
        this->dbConn.get(), query.c_str(), -1, &stmt, nullptr);

    if (ret != SQLITE_OK) {
        this->error = sqlite3_errstr(ret);
        return std::nullopt;
    }

    std::vector<model::Binary> binList;
    while ((ret = sqlite3_step(stmt)) == SQLITE_ROW) {
        auto name =
            reinterpret_cast<const char *>(sqlite3_column_text(stmt, 1));
        auto path =
            reinterpret_cast<const char *>(sqlite3_column_text(stmt, 2));

        binList.emplace_back(name, path);
    }

    sqlite3_finalize(stmt);
    if (ret != SQLITE_DONE) {
        this->error = sqlite3_errstr(ret);
        return {};
    }

    return binList;
}

std::string BinaryDao::getError() {
    return this->error;
}
};  // namespace OcelotMDM::component::db
