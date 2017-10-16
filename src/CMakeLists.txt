
set(ROOT_FOLDER "plugins")

function(list_subdirectories)
	set(options PARENT)
    set(singleValueOpts VAR DIR)

    cmake_parse_arguments(LIST_SUBDIR "${options}" "${singleValueOpts}" "${multiValueOpts}" ${ARGN})

    if (NOT DEFINED LIST_SUBDIR_VAR)
    	message(FATAL_ERROR "VAR not specified for list_subdirectories")
    endif()

    if (NOT DEFINED LIST_SUBDIR_DIR)
    	message(FATAL_ERROR "DIR not specified for list_subdirectories")
    endif()

   file(GLOB sub-dir ${LIST_SUBDIR_DIR}/*)

   set(list_of_dirs "")

   foreach(dir ${sub-dir})
     if(IS_DIRECTORY ${curdir}/${dir})
         set(list_of_dirs ${list_of_dirs} ${dir})
     endif()
   endforeach()

   set(${LIST_SUBDIR_VAR} ${list_of_dirs} PARENT_SCOPE)

endfunction()

function(add_plugin)
    set(options)
    set(singleValueOpts NAME ROOT MANIFEST EXECUTABLE)

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


    if (NOT DEFINED ADD_PLUGIN_MANIFEST)
        set(ADD_PLUGIN_MANIFEST "${CMAKE_CURRENT_SOURCE_DIR}/manifest.json")
    endif()

    if (NOT DEFINED ADD_PLUGIN_EXECUTABLE)
        set(ADD_PLUGIN_EXECUTABLE "${CMAKE_CURRENT_BINARY_DIR}/${ADD_PLUGIN_NAME}")
    endif()

    set(OUTPUT_FOLDER "${CMAKE_BINARY_DIR}/gen/${ADD_PLUGIN_ROOT}/${ADD_PLUGIN_NAME}")

    file(MAKE_DIRECTORY ${OUTPUT_FOLDER})

    if (EXISTS "${ADD_PLUGIN_MANIFEST}.in")
      get_filename_component(MANIFEST_BASENAME ${ADD_PLUGIN_MANIFEST} NAME)
      configure_file("${ADD_PLUGIN_MANIFEST}.in" "${OUTPUT_FOLDER}/${MANIFEST_BASENAME}" @ONLY)
      add_custom_command(TARGET ${ADD_PLUGIN_NAME} POST_BUILD COMMAND ${CMAKE_COMMAND} ARGS -E copy ${ADD_PLUGIN_EXECUTABLE} ${OUTPUT_FOLDER})
    else()
      add_custom_command(TARGET ${ADD_PLUGIN_NAME} POST_BUILD COMMAND ${CMAKE_COMMAND} ARGS -E copy ${ADD_PLUGIN_MANIFEST} ${ADD_PLUGIN_EXECUTABLE} ${OUTPUT_FOLDER})
    endif()
endfunction()

list_subdirectories(VAR PREP_PLUGINS DIR ${CMAKE_CURRENT_LIST_DIR})

foreach(PLUGIN ${PREP_PLUGINS})
	add_subdirectory(${PLUGIN})
endforeach()

