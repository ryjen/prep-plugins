# set variables
set(SOURCE_FILES main.go)

set(BINARY_NAME main)

get_filename_component(PLUGIN_NAME ${CMAKE_CURRENT_LIST_DIR} NAME)

# setup plugin
add_go_executable(TARGET ${PLUGIN_NAME} NAME ${BINARY_NAME} SOURCE_FILES ${SOURCE_FILES})

add_plugin(TARGET ${PLUGIN_NAME} EXECUTABLE ${BINARY_NAME})

# setup test
configure_file(main_test.go.in ${CMAKE_CURRENT_LIST_DIR}/main_test.go)

set(TEST_NAME ${PLUGIN_NAME}-test)

add_test(NAME ${TEST_NAME}
			WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}
            COMMAND ${GOSHIM} test -data ${TEST_DATA_PATH})
