
set(SOURCE_FILES main.go)

get_filename_component(PLUGIN_NAME "${CMAKE_CURRENT_LIST_DIR}" NAME)

add_executable(${PLUGIN_NAME} ${SOURCE_FILES})

add_plugin(NAME ${PLUGIN_NAME})
