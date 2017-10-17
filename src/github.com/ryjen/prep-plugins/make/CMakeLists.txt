
set(SOURCE_FILES main.go)

add_go_executable(NAME make SOURCE_FILES ${SOURCE_FILES})

add_plugin(NAME make)
