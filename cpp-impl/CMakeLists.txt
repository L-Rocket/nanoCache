cmake_minimum_required(VERSION 3.10)
project(main)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include_directories(include)

enable_testing()
add_executable(test_SimpleCache tests/test_SimpleCache.cpp src/SimpleCache.cpp)
add_test(NAME test_SimpleCache COMMAND test_SimpleCache)

add_executable(main src/main.cpp src/SimpleCache.cpp)

add_executable(cli apps/cli.cpp src/SimpleCache.cpp)
