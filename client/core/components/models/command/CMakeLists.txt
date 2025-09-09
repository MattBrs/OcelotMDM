set(LIBRARY_NAME ocelot_model_command)

add_library(${LIBRARY_NAME} INTERFACE)

target_include_directories(
  ${LIBRARY_NAME}
  PUBLIC
  INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})

add_library(OcelotMDM::Component::Model::Command ALIAS ${LIBRARY_NAME})
