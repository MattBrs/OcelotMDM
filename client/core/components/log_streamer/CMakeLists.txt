set(LIBRARY_NAME ocelot_log_streamer)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/log_streamer.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/log_streamer.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME} PUBLIC OcelotMDM::Component::Network::MqttClient
                         OcelotMDM::Component::Timer)

add_library(OcelotMDM::Component::LogStreamer ALIAS ${LIBRARY_NAME})
