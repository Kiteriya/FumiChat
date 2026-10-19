#ifndef FUMICLIENT_CONNECTIONHANDLER_H
#define FUMICLIENT_CONNECTIONHANDLER_H

#include <QtNetwork/QNetworkAccessManager>
#include <QtNetwork/QNetworkRequest>
#include <QtNetwork/QNetworkReply>
#include <QtWebSockets/QtWebSockets>
#include "DTO/userAuth_DTO.h"

class ConnectionHandler : public QObject{
    Q_OBJECT
public:
    explicit ConnectionHandler(QObject *parent = nullptr);

    void userRegister(const UserRegistration &data);
    void userLogin(const UserRegistration &data);
    void userConnectionWebsocket();

private:
    QNetworkAccessManager m_networkManager;
    QString m_baseURL;
    QString token;

    QNetworkRequest createAuthRequest(const QString &endpoint);
    QNetworkRequest createUnauthRequest(const QString &endpoint);
    
};

#endif //FUMICLIENT_CONNECTIONHANDLER_H
