
cmake_minimum_required (VERSION 3.0)
project (epldbg)

add_compile_options(-std=c++14)

include_directories(libelfin linenoise )
add_executable(epldbg luncher.cpp linenoise.c)

add_custom_target(
   libelfin
   COMMAND make
   WORKING_DIRECTORY ${PROJECT_SOURCE_DIR}/libelfin
)
target_link_libraries(epldbg
                      ${PROJECT_SOURCE_DIR}/libelfin/dwarf/libdwarf++.so
                      ${PROJECT_SOURCE_DIR}/libelfin/elf/libelf++.so)
add_dependencies(epldbg libelfin)
