set(LIBRARY_NAME ocelot_api_dto)

add_library(${LIBRARY_NAME} INTERFACE)

target_include_directories(
  ${LIBRARY_NAME}
  PUBLIC
  INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})

add_library(OcelotMDM::Component::Api::Dto ALIAS ${LIBRARY_NAME})
