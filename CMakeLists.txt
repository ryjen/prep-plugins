cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins LANGUAGES Go VERSION 0.1)

set(GOPATH "${CMAKE_CURRENT_SOURCE_DIR}")
file(MAKE_DIRECTORY ${GOPATH})

include(golang)

file(GLOB PLUGIN_DIRS src/github.com/ryjen/prep-plugins/*)

set(PLUGIN_NAMES "")

foreach(PLUGIN_DIR ${PLUGIN_DIRS})
	if(IS_DIRECTORY ${PLUGIN_DIR})
		get_filename_component(PLUGIN_NAME ${PLUGIN_DIR} NAME)
		list(APPEND PLUGIN_NAMES ${PLUGIN_NAME})
	endif()
endforeach()

list(REMOVE_ITEM PLUGIN_NAMES support)

list(INSERT PLUGIN_NAMES 0 support)

set(PREP_PLUGINS ${PLUGIN_NAMES} PARENT_SCOPE)

add_subdirectory(src)

