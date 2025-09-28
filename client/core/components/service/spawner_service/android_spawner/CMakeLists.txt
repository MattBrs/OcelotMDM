set(LIBRARY_NAME ocelot_android_spawner_service)

add_library(
  ${LIBRARY_NAME} STATIC
  ${CMAKE_CURRENT_SOURCE_DIR}/android_spawner_service.hpp
  ${CMAKE_CURRENT_SOURCE_DIR}/android_spawner_service.cpp)

target_include_directories(${LIBRARY_NAME} PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})
target_link_libraries(
  ${LIBRARY_NAME} PUBLIC OcelotMDM::Component::Logger
                         OcelotMDM::Component::Service::SpawnerService)

add_library(OcelotMDM::Component::Service::AndroidSpawnerService ALIAS
            ${LIBRARY_NAME})
