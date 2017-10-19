cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins LANGUAGES Go VERSION 0.1)

set(GOPATH "${CMAKE_CURRENT_SOURCE_DIR}")

include(golang)

add_subdirectory(src)

add_custom_target(build_plugins ALL DEPENDS support make cmake autotools git homebrew archive)