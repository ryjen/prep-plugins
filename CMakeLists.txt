cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins LANGUAGES Go VERSION 0.1)

set(GOPATH "${CMAKE_CURRENT_SOURCE_DIR}")

include(golang)

enable_testing()

add_subdirectory(src)

set(PLUGINS support make cmake autotools git homebrew archive)

add_custom_target(build_plugins ALL DEPENDS ${PLUGINS})
