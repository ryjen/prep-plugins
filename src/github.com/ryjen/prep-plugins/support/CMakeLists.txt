
get_filename_component(PLUGIN_NAME ${CMAKE_CURRENT_LIST_DIR} NAME)

add_go_library(TARGET ${PLUGIN_NAME} NAME plugin SOURCE_FILES plugin.go)

add_plugin(TARGET ${PLUGIN_NAME} EXECUTABLE plugin)
