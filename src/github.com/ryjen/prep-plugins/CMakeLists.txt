
set(ROOT_FOLDER "plugins")

function(add_plugin)
    set(options)
    set(singleValueOpts TARGET ROOT EXECUTABLE)

    cmake_parse_arguments(ADD_PLUGIN "${options}" "${singleValueOpts}" "${multiValueOpts}" ${ARGN})

    if (NOT DEFINED ADD_PLUGIN_TARGET)
        message(FATAL_ERROR "Target for add_plugin not specified")
    endif()

    if (NOT DEFINED ADD_PLUGIN_ROOT)
        set(ADD_PLUGIN_ROOT ${ROOT_FOLDER})
    endif()

    if (NOT DEFINED ADD_PLUGIN_EXECUTABLE)
        set(ADD_PLUGIN_EXECUTABLE "${CMAKE_CURRENT_BINARY_DIR}/${ADD_PLUGIN_TARGET}")
    endif()

    set(MANIFEST_NAME "manifest.json")

    set(OUTPUT_FOLDER "${CMAKE_BINARY_DIR}/gen/${ADD_PLUGIN_ROOT}/${ADD_PLUGIN_TARGET}")

    file(MAKE_DIRECTORY ${OUTPUT_FOLDER})

    if (EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/${MANIFEST_NAME}.in")
      configure_file("${CMAKE_CURRENT_SOURCE_DIR}/${MANIFEST_NAME}.in" "${OUTPUT_FOLDER}/${MANIFEST_NAME}" @ONLY)
    endif()

    add_custom_command(TARGET ${ADD_PLUGIN_TARGET} POST_BUILD COMMAND ${CMAKE_COMMAND} ARGS -E copy ${ADD_PLUGIN_EXECUTABLE}
            ${OUTPUT_FOLDER})

endfunction()

# standard plugins

add_subdirectory(support)

add_subdirectory(make)

add_subdirectory(cmake)

add_subdirectory(autotools)

add_subdirectory(git)

add_subdirectory(homebrew)

add_subdirectory(archive)

# plugin dependencies

add_dependencies(make support)

add_dependencies(cmake support)

add_dependencies(autotools support)

add_dependencies(git support)

add_dependencies(homebrew support)

add_dependencies(archive support)
