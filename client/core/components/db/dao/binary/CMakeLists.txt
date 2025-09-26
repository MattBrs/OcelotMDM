set(LIBRARY_NAME ocelot_component_dao_binary)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/binary_dao.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/binary_dao.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME}
  PUBLIC SQLite::SQLite3 OcelotMDM::Component::Db::Migrations
         OcelotMDM::Component::Db::Shared OcelotMDM::Component::Model::Binary)

add_library(OcelotMDM::Component::Db::Dao::Binary ALIAS ${LIBRARY_NAME})
