
add_go_library(NAME support SOURCE_FILES plugin.go)

add_plugin(NAME support EXECUTABLE libsupport.a)
