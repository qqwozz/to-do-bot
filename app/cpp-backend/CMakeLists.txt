cmake_minimum_required(VERSION 3.14)
project(todo-backend VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(SQLITE3 REQUIRED sqlite3)
find_package(OpenSSL REQUIRED)

if(NOT CROW_INCLUDE_DIR)
    find_path(CROW_INCLUDE_DIR crow.h
        HINTS /usr/local/include
        HINTS /usr/include
    )
endif()

if(NOT CROW_INCLUDE_DIR)
    message(FATAL_ERROR "Crow headers not found. Set CROW_INCLUDE_DIR or install Crow.")
endif()

set(SOURCES
    main.cpp
    src/database.cpp
    src/handlers.cpp
)

set(HEADERS
    include/models.h
    include/database.h
    include/handlers.h
)

add_executable(${PROJECT_NAME} ${SOURCES} ${HEADERS})

target_include_directories(${PROJECT_NAME} PRIVATE
    ${CMAKE_SOURCE_DIR}/include
    ${SQLITE3_INCLUDE_DIRS}
    ${CROW_INCLUDE_DIR}
)

target_link_libraries(${PROJECT_NAME} PRIVATE
    ${SQLITE3_LIBRARIES}
    OpenSSL::SSL
    OpenSSL::Crypto
    pthread
)

option(BUILD_TESTS "Build tests" OFF)
install(TARGETS ${PROJECT_NAME} DESTINATION bin)
