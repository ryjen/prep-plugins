cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins LANGUAGES Go VERSION 0.1)

set(GOPATH "${CMAKE_CURRENT_SOURCE_DIR}")
file(MAKE_DIRECTORY ${GOPATH})

include(golang)

set(PREP_PLUGINS support make PARENT_SCOPE)

add_subdirectory(src)

