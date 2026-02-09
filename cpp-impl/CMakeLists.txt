cmake_minimum_required(VERSION 3.10)
project(main)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include_directories(include)

enable_testing()
add_executable(test_SimpleCache tests/test_SimpleCache.cpp src/SimpleCache.cpp)
add_test(NAME test_SimpleCache COMMAND test_SimpleCache)

add_executable(test_ShardedCache_basic tests/test_ShardedCache_basic.cpp src/ShardedCache.cpp src/CacheShared.cpp)
add_test(NAME test_ShardedCache_basic COMMAND test_ShardedCache_basic)

add_executable(test_ShardedCache_delete tests/test_ShardedCache_delete.cpp src/ShardedCache.cpp src/CacheShared.cpp)
add_test(NAME test_ShardedCache_delete COMMAND test_ShardedCache_delete)

add_executable(test_ShardedCache_ttl tests/test_ShardedCache_ttl.cpp src/ShardedCache.cpp src/CacheShared.cpp)
add_test(NAME test_ShardedCache_ttl COMMAND test_ShardedCache_ttl)

add_executable(test_ShardedCache_concurrent tests/test_ShardedCache_concurrent.cpp src/ShardedCache.cpp src/CacheShared.cpp)
add_test(NAME test_ShardedCache_concurrent COMMAND test_ShardedCache_concurrent)

add_executable(main src/main.cpp src/SimpleCache.cpp)

add_executable(cli apps/cli.cpp src/SimpleCache.cpp)
