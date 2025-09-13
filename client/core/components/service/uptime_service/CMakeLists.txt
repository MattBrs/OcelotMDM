set(LIBRARY_NAME ocelot_uptime_service)

add_library(
  ${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/uptime_service.hpp
                         ${CMAKE_CURRENT_SOURCE_DIR}/uptime_service.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(${LIBRARY_NAME}
                      PUBLIC OcelotMDM::Component::Network::MqttClient)

add_library(OcelotMDM::Component::Service::UptimeService ALIAS ${LIBRARY_NAME})
