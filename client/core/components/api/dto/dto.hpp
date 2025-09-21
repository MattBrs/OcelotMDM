#pragma once

#include <string>

namespace OcelotMDM::component::Api::Dto {
typedef struct getBinaryReq {
    std::string binaryName;
    std::string otp;
} GetBinaryReq;

typedef struct getBinaryRes {
    std::string name;
    std::string version;
    std::string binaryData;
} GetBinaryRes;
}  // namespace OcelotMDM::component::Api::Dto
