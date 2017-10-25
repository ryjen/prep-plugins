
# set variables
set(SOURCE_FILES main.go)

set(BINARY_NAME main)

# setup plugin
external_go_project(TARGET goarchive PACKAGE github.com/mholt/archiver)

get_filename_component(PLUGIN_NAME ${CMAKE_CURRENT_LIST_DIR} NAME)

add_go_executable(TARGET ${PLUGIN_NAME} NAME ${BINARY_NAME} SOURCE_FILES ${SOURCE_FILES})

add_dependencies(${PLUGIN_NAME} goarchive)

add_plugin(TARGET ${PLUGIN_NAME} EXECUTABLE ${BINARY_NAME})

# setup test
set(TEST_SOURCES main_test.go)

set(TEST_NAME ${PLUGIN_NAME}-test)

add_go_executable(TARGET ${TEST_NAME} NAME ${TEST_NAME} SOURCE_FILES ${SOURCE_FILES} ${TEST_SOURCES})

add_test(TARGET ${TEST_NAME} EXECUTABLE ${CMAKE_CURRENT_BINARY_DIR}/${TEST_NAME})
