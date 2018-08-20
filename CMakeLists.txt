cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins VERSION 0.1)

enable_testing()

set(PLUGIN_BINARY "main")

add_subdirectory(src)
add_subdirectory(tests)

install(DIRECTORY ${CMAKE_BINARY_DIR}/gen/plugins DESTINATION ${CMAKE_INSTALL_PREFIX}/share/prep PATTERN ${PLUGIN_BINARY} PERMISSIONS OWNER_WRITE OWNER_READ OWNER_EXECUTE GROUP_READ GROUP_EXECUTE WORLD_READ WORLD_EXECUTE)


