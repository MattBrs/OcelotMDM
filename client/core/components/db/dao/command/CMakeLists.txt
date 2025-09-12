set(LIBRARY_NAME ocelot_component_dao_command)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/command_dao.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/command_dao.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME}
  PUBLIC SQLite::SQLite3 OcelotMDM::Component::Db::Migrations
         OcelotMDM::Component::Model::Command OcelotMDM::Component::Db::Shared)

add_library(OcelotMDM::Component::Db::Dao::Command ALIAS ${LIBRARY_NAME})
