cmake_minimum_required(VERSION 3.25)
project(ocelot_mdm_core CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_package(CURL REQUIRED)
find_package(PahoMqttCpp REQUIRED)
find_package(msgpack-cxx REQUIRED)

add_subdirectory(${CMAKE_SOURCE_DIR}/components/network)

add_executable(${PROJECT_NAME} main.cpp)

target_link_libraries(
  ${PROJECT_NAME} CURL::libcurl OcelotMDM::Component::Network::HttpClient
  OcelotMDM::Component::Network::MqttClient msgpack-cxx)
