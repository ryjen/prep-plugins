
add_go_library(NAME ${PLUGIN} SOURCE_FILES plugin.go)

add_plugin(NAME ${PLUGIN} EXECUTABLE lib${PLUGIN}.a)
