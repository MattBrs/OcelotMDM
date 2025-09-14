cmake_minimum_required(VERSION 3.25)
project(ocelot_mdm_core CXX)

if(BUILD_TYPE STREQUAL "Release")
  set(OPTIM "-O2 -pthread")
  message("Building for release")
else()
  set(OPTIM "-O0 -g -Wall -Wextra -pthread")
  message("Building for debug")
endif()

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} ${OPTIM} -fdiagnostics-color=always")

find_package(SQLite3 REQUIRED)
find_package(CURL REQUIRED)
find_package(nlohmann_json REQUIRED)

add_subdirectory(${CMAKE_SOURCE_DIR}/components/db)
add_subdirectory(${CMAKE_SOURCE_DIR}/components/logger)
add_subdirectory(${CMAKE_SOURCE_DIR}/components/models)
add_subdirectory(${CMAKE_SOURCE_DIR}/components/network)
add_subdirectory(${CMAKE_SOURCE_DIR}/components/service)
add_subdirectory(${CMAKE_SOURCE_DIR}/components/commands_impl)

add_executable(${PROJECT_NAME} main.cpp)

target_link_libraries(
  ${PROJECT_NAME}
  OcelotMDM::Component::Network::HttpClient
  OcelotMDM::Component::Network::MqttClient
  OcelotMDM::Component::Service::CommandService
  OcelotMDM::Component::Service::UptimeService
  OcelotMDM::Component::Db::Client)
