set(LIBRARY_NAME ocelot_utils)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/utils.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/utils.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(${LIBRARY_NAME} PUBLIC OcelotMDM::Component::Logger)

add_library(OcelotMDM::Component::Utils ALIAS ${LIBRARY_NAME})
