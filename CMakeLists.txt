cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins LANGUAGES Go VERSION 0.1)

set(GOPATH "${CMAKE_CURRENT_SOURCE_DIR}")

include(golang)

enable_testing()

set(PLUGINS support make cmake autotools git homebrew archive)

add_custom_target(build_plugins ALL DEPENDS ${PLUGINS})

set(GOSHIM ${CMAKE_BINARY_DIR}/go_shim)

configure_file(${CMAKE_CURRENT_LIST_DIR}/go_shim.in ${GOSHIM} @ONLY)

set(TEST_DATA_PATH ${CMAKE_CURRENT_LIST_DIR}/tests/data)

add_subdirectory(src)
