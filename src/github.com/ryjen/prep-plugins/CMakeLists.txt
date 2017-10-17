
set(ROOT_FOLDER "plugins")

function(add_plugin)
    set(options)
    set(singleValueOpts NAME ROOT EXECUTABLE)

    cmake_parse_arguments(ADD_PLUGIN "${options}" "${singleValueOpts}" "${multiValueOpts}" ${ARGN})

    if (NOT DEFINED ADD_PLUGIN_NAME)
        message(FATAL_ERROR "FOLDER name for add_plugin not specified")
    endif()

    if (NOT DEFINED ADD_PLUGIN_ROOT)
        if (NOT DEFINED ROOT_FOLDER)
            message(FATAL_ERROR "ROOT folder for add_plugin not specified")
        endif()
        set(ADD_PLUGIN_ROOT ${ROOT_FOLDER})
    endif()

    set(MANIFEST_NAME "manifest.json")

    if (NOT DEFINED ADD_PLUGIN_EXECUTABLE)
        set(ADD_PLUGIN_EXECUTABLE "${CMAKE_CURRENT_BINARY_DIR}/${ADD_PLUGIN_NAME}")
    endif()

    set(OUTPUT_FOLDER "${CMAKE_BINARY_DIR}/gen/${ADD_PLUGIN_ROOT}/${ADD_PLUGIN_NAME}")

    file(MAKE_DIRECTORY ${OUTPUT_FOLDER})

    if (EXISTS "${CMAKE_CURRENT_SOURCE_DIR}/${MANIFEST_NAME}.in")
      configure_file("${CMAKE_CURRENT_SOURCE_DIR}/${MANIFEST_NAME}.in" "${OUTPUT_FOLDER}/${MANIFEST_NAME}" @ONLY)
    endif()

    add_custom_command(TARGET ${ADD_PLUGIN_NAME} POST_BUILD COMMAND ${CMAKE_COMMAND} ARGS -E copy ${ADD_PLUGIN_EXECUTABLE} ${OUTPUT_FOLDER})

endfunction()

foreach(PLUGIN ${PLUGIN_NAMES})
  message("Adding plugin ${PLUGIN}")
  add_subdirectory(${PLUGIN})
endforeach()


add_dependencies(make support)

