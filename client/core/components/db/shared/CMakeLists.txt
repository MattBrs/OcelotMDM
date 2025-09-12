set(LIBRARY_NAME ocelot_component_db_shared)

add_library(${LIBRARY_NAME} INTERFACE)

target_include_directories(
  ${LIBRARY_NAME}
  PUBLIC
  INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})

add_library(OcelotMDM::Component::Db::Shared ALIAS ${LIBRARY_NAME})
