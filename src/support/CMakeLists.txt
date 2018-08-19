
set(BINARY_NAME plugin)

get_filename_component(PLUGIN_NAME ${CMAKE_CURRENT_LIST_DIR} NAME)

add_custom_target(${PLUGIN_NAME})

add_plugin(TARGET ${PLUGIN_NAME} EXECUTABLE ${BINARY_NAME})


