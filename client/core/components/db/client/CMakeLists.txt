set(LIBRARY_NAME ocelot_component_db_client)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/db_client.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/db_client.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME} PUBLIC SQLite::SQLite3 OcelotMDM::Component::Db::Migrations
                         OcelotMDM::Component::Db::Dao::Command)

add_library(OcelotMDM::Component::Db::Client ALIAS ${LIBRARY_NAME})
