cmake_minimum_required(VERSION 3.14)
project(todo-backend-tests LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(SQLITE3 REQUIRED sqlite3)
find_package(GTest REQUIRED)
find_package(OpenSSL REQUIRED)

if(NOT CROW_INCLUDE_DIR)
    find_path(CROW_INCLUDE_DIR crow.h
        HINTS /usr/local/include
        HINTS /usr/include
    )
endif()

include_directories(${CMAKE_SOURCE_DIR}/../include)
include_directories(${GTEST_INCLUDE_DIRS})

if(CROW_INCLUDE_DIR)
    include_directories(${CROW_INCLUDE_DIR})
endif()

set(DATABASE_SOURCES
    ${CMAKE_SOURCE_DIR}/../src/database.cpp
)

add_executable(database_tests database_test.cpp ${DATABASE_SOURCES})
target_include_directories(database_tests PRIVATE
    ${CMAKE_SOURCE_DIR}/../include
    ${SQLITE3_INCLUDE_DIRS}
    ${CROW_INCLUDE_DIR}
)
target_link_libraries(database_tests
    GTest::GTest
    GTest::Main
    ${SQLITE3_LIBRARIES}
    pthread
)

add_executable(model_tests model_test.cpp)
target_include_directories(model_tests PRIVATE
    ${CMAKE_SOURCE_DIR}/../include
    ${CROW_INCLUDE_DIR}
)
target_link_libraries(model_tests
    GTest::GTest
    GTest::Main
    pthread
)

enable_testing()
add_test(NAME DatabaseTests COMMAND database_tests)
add_test(NAME ModelTests COMMAND model_tests)

if(CROW_INCLUDE_DIR)
    set(HANDLER_SOURCES
        ${CMAKE_SOURCE_DIR}/../src/database.cpp
        ${CMAKE_SOURCE_DIR}/../src/handlers.cpp
    )

    add_executable(handler_tests handler_test.cpp ${HANDLER_SOURCES})
    target_include_directories(handler_tests PRIVATE
        ${CMAKE_SOURCE_DIR}/../include
        ${SQLITE3_INCLUDE_DIRS}
        ${CROW_INCLUDE_DIR}
    )
    target_link_libraries(handler_tests
        GTest::GTest
        GTest::Main
        ${SQLITE3_LIBRARIES}
        OpenSSL::SSL
        OpenSSL::Crypto
        pthread
    )
    add_test(NAME HandlerTests COMMAND handler_tests)
else()
    message(STATUS "Crow not found, skipping handler tests")
endif()
