cmake_minimum_required(VERSION 4.3)

project(FumiClient LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Quick Network)
qt_policy(SET QTP0001 NEW)
qt_policy(SET QTP0004 NEW)

qt_standard_project_setup()

qt_add_executable(FumiClient
        src/main.cpp
        src/networkManager/connectionHandler.h
        src/networkManager/connectionHandler.h
        src/networkManager/connectionHandler.cpp
        src/networkManager/DTO/userAuth_DTO.h
)

qt_add_qml_module(FumiClient
        URI FumiClient
        VERSION 1.0
        QML_FILES
        ui/Main.qml
)

target_link_libraries(FumiClient
        PRIVATE Qt6::Quick
        PRIVATE Qt6::Network
)