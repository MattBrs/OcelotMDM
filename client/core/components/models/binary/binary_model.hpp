#pragma once

#include <string>

namespace OcelotMDM::component::model {
class Binary {
   public:
    Binary(const std::string &name, const std::string &path);

    std::string getId() const;
    std::string getName() const;
    std::string getPath() const;

    void setId(const std::string &id);
    void setName(const std::string &name);
    void setPath(const std::string &path);

   private:
    std::string id;
    std::string name;
    std::string path;
};
}  // namespace OcelotMDM::component::model
