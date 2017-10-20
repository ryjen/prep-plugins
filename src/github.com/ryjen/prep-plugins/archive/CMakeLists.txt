
set(SOURCE_FILES main.go)

set(BINARY_NAME main)

external_go_project(TARGET goarchive PACKAGE github.com/mholt/archiver)

get_filename_component(PLUGIN_NAME ${CMAKE_CURRENT_LIST_DIR} NAME)

add_go_executable(TARGET ${PLUGIN_NAME} NAME ${BINARY_NAME} SOURCE_FILES ${SOURCE_FILES})

add_go_test(TARGET ${PLUGIN_NAME}_test SOURCE_FILES ${SOURCE_FILES} archive_test.go)

add_dependencies(${PLUGIN_NAME} goarchive)

add_plugin(TARGET ${PLUGIN_NAME} EXECUTABLE ${BINARY_NAME})
