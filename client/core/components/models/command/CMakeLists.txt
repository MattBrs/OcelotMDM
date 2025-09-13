set(LIBRARY_NAME ocelot_model_command)

add_library(
  ${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/command_model.hpp
                         ${CMAKE_CURRENT_SOURCE_DIR}/command_model.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

add_library(OcelotMDM::Component::Model::Command ALIAS ${LIBRARY_NAME})
