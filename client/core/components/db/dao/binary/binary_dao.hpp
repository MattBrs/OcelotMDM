#pragma once

#include <sqlite3.h>

#include <memory>
#include <optional>
#include <string>
#include <vector>
namespace OcelotMDM::component::db {
class BinaryDao {
   public:
    explicit BinaryDao(const std::shared_ptr<sqlite3> &dbConn);

    std::optional<bool> addBinary(
        const std::string &name, const std::string &path);
    std::optional<bool> removeBinary(const std::string &name);
    std::optional<std::vector<std::string>> listBinaries();

    std::string getError();

   private:
    std::string              error;
    std::shared_ptr<sqlite3> dbConn;
};
};  // namespace OcelotMDM::component::db
