cmake_minimum_required(VERSION 3.14)
project(todo-backend VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(SQLITE3 REQUIRED sqlite3)
find_package(Crow REQUIRED)

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
)

target_link_libraries(${PROJECT_NAME} PRIVATE
    Crow::Crow
    ${SQLITE3_LIBRARIES}
    pthread
)

option(BUILD_TESTS "Build tests" OFF)
install(TARGETS ${PROJECT_NAME} DESTINATION bin)
