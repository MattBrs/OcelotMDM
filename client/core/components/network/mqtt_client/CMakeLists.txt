set(LIBRARY_NAME ocelot_mqtt_client)

find_package(PahoMqttCpp REQUIRED)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/mqtt_client.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/mqtt_client.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(${LIBRARY_NAME} PUBLIC PahoMqttCpp::paho-mqttpp3-static)

add_library(OcelotMDM::Component::Network::MqttClient ALIAS ${LIBRARY_NAME})
