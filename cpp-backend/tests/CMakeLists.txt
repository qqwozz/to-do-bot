cmake_minimum_required(VERSION 3.14)
project(todo-backend-tests LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(SQLITE3 REQUIRED sqlite3)
find_package(Crow REQUIRED)
find_package(GTest REQUIRED)

include_directories(${CMAKE_SOURCE_DIR}/../include)
include_directories(${GTEST_INCLUDE_DIRS})

set(SOURCES
    ${CMAKE_SOURCE_DIR}/../src/database.cpp
)

add_executable(database_tests database_test.cpp ${SOURCES})
target_include_directories(database_tests PRIVATE
    ${CMAKE_SOURCE_DIR}/../include
    ${SQLITE3_INCLUDE_DIRS}
)
target_link_libraries(database_tests
    GTest::GTest
    GTest::Main
    ${SQLITE3_LIBRARIES}
    pthread
)

enable_testing()
add_test(NAME DatabaseTests COMMAND database_tests)
