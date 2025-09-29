#include "linux_spawner_service.hpp"

#include <sched.h>
#include <unistd.h>

#include <cstdlib>
#include <iostream>
#include <memory>

#include "binary_dao.hpp"
#include "spawner_service.hpp"

namespace OcelotMDM::component::service {
LinuxSpawerService::LinuxSpawerService(
    const std::shared_ptr<db::BinaryDao> &binDao)
    : SpawnerService(binDao) {
    this->startBinaries();
}

int LinuxSpawerService::runBinary(const std::string &binPath) {
    std::cout << "about to run binary: " << binPath << std::endl;
    pid_t pid = fork();
    if (pid == 0) {
        char *const args[] = {const_cast<char *>(binPath.c_str()), nullptr};
        execvp(binPath.c_str(), args);
        exit(1);  // child failed the execution
    }

    if (pid < 0) {
        return -1;
    }

    return pid;
}
}  // namespace OcelotMDM::component::service
