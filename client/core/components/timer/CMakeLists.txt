set(LIBRARY_NAME ocelot_timer)

add_library(${LIBRARY_NAME} STATIC ${CMAKE_CURRENT_SOURCE_DIR}/timer.hpp
                                   ${CMAKE_CURRENT_SOURCE_DIR}/timer.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
# target_link_libraries(${LIBRARY_NAME} PUBLIC OcelotMDM::Component::Logger)

add_library(OcelotMDM::Component::Timer ALIAS ${LIBRARY_NAME})
