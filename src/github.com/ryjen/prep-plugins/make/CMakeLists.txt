
set(SOURCE_FILES main.go)

add_go_executable(NAME ${PLUGIN} SOURCE_FILES ${SOURCE_FILES})

add_plugin(NAME ${PLUGIN})
