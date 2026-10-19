#include <iostream>
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include "networkManager/connectionHandler.h"


int main(int argc, char* argv[]) {

    //Qml application
    QGuiApplication app(argc, argv);
    QQmlApplicationEngine engine;
    engine.loadFromModule("FumiClient", "Main");

    //Network connection
    auto *connectionHandler = new ConnectionHandler(&app);


    UserRegistration regData;
    regData.username = "Yuki";
    regData.password = "boba123";
     connectionHandler->userRegister(regData);

    connectionHandler->userLogin(regData);

    return app.exec();
}
