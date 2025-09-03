set(LIBRARY_NAME ocelot_http_client)

add_library(
  ${LIBRARY_NAME} STATIC
  ${CMAKE_CURRENT_SOURCE_DIR}/http_client.hpp
  ${CMAKE_CURRENT_SOURCE_DIR}/http_client.cpp
  ${CMAKE_CURRENT_SOURCE_DIR}/http_response.hpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(${LIBRARY_NAME} PUBLIC CURL::libcurl)

add_library(OcelotMDM::Component::Network::HttpClient ALIAS ${LIBRARY_NAME})
