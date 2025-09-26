set(LIBRARY_NAME ocelot_commands_impl)

add_library(
  ${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/commands_impl.hpp
                         ${CMAKE_CURRENT_SOURCE_DIR}/commands_impl.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME}
  PUBLIC OcelotMDM::Component::Network::MqttClient
         OcelotMDM::Component::Network::HttpClient
         OcelotMDM::Component::Logger
         OcelotMDM::Component::Utils
         OcelotMDM::Component::Api::Handler
         OcelotMDM::Component::Utils
         OcelotMDM::Component::Timer
         OcelotMDM::Component::LogStreamer
         OcelotMDM::Component::Db::Dao::Binary)

add_library(OcelotMDM::Component::CommandsImpl ALIAS ${LIBRARY_NAME})
