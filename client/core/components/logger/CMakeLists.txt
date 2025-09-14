set(LIBRARY_NAME ocelot_logger)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/logger.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/logger.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

add_library(OcelotMDM::Component::Logger ALIAS ${LIBRARY_NAME})
