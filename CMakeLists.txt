cmake_minimum_required (VERSION 3.1)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake")

# define project name
project (prep-plugins VERSION 0.1)

enable_testing()

add_subdirectory(src)
add_subdirectory(tests)

