cmake_minimum_required (VERSION 3.8)
project ("DFA-Toolkit")
set(CMAKE_CXX_STANDARD 17)

set(MODULE_SOURCE "${PROJECT_SOURCE_DIR}/source/module")
set(EXTERNALS "${PROJECT_SOURCE_DIR}/externals")

# pybind11
add_subdirectory(${EXTERNALS}/pybind11)

# module
pybind11_add_module (
    DFA_Toolkit 
    ${MODULE_SOURCE}/DFA.cpp
    ${MODULE_SOURCE}/DFA.h
)

target_include_directories (
    DFA_Toolkit 
    PRIVATE ${MODULE_SOURCE}
)