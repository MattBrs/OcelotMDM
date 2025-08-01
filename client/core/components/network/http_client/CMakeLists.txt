set(LIBRARY_NAME ocelot_http_client)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/http_client.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/http_client.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

add_library(OcelotMDM::Component::HttpClient ALIAS ${LIBRARY_NAME})
