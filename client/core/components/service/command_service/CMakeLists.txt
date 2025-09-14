set(LIBRARY_NAME ocelot_command_service)

add_library(
  ${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/command_service.hpp
                         ${CMAKE_CURRENT_SOURCE_DIR}/command_service.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME}
  PUBLIC OcelotMDM::Component::Network::MqttClient
         OcelotMDM::Component::Model::Command
         nlohmann_json::nlohmann_json
         OcelotMDM::Component::Db::Dao::Command
         OcelotMDM::Component::CommandsImpl)

add_library(OcelotMDM::Component::Service::CommandService ALIAS ${LIBRARY_NAME})
