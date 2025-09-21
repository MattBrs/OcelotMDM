set(LIBRARY_NAME ocelot_api_handler)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/api_handler.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/api_handler.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME}
  PUBLIC OcelotMDM::Component::Network::HttpClient OcelotMDM::Component::Logger
         OcelotMDM::Component::Utils OcelotMDM::Component::Api::Dto
         nlohmann_json::nlohmann_json)

add_library(OcelotMDM::Component::Api::Handler ALIAS ${LIBRARY_NAME})
