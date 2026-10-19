#include "connectionHandler.h"

ConnectionHandler::ConnectionHandler(QObject *parent)
    : QObject(parent),
    m_networkManager(this),
    m_baseURL("http://localhost:8080///FumiServer/")
    {}

void ConnectionHandler::userRegister(const UserRegistration &data) {
    const QNetworkRequest request = createUnauthRequest("register");

    QJsonObject json_object;

    json_object["username"] = data.username;
    json_object["password"] = data.password;

    const QByteArray post_data = QJsonDocument(json_object).toJson(QJsonDocument::Compact);

    QNetworkReply *reply = m_networkManager.post(request,post_data);

    QObject::connect(reply,&QNetworkReply::finished,[reply]() {
        reply->deleteLater();

        if (reply->error() != QNetworkReply::NoError) {
            qWarning() <<"Network Error"<< reply-> errorString();
            return;
        }

        if (int StatusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();
            StatusCode >=200 && StatusCode <300){
        }

    });
}

void ConnectionHandler::userLogin(const UserRegistration &data) {
    const QNetworkRequest request = createUnauthRequest("login");

    QJsonObject json_object;

    json_object["username"] = data.username;
    json_object["password"] = data.password;

    const QByteArray post_data = QJsonDocument(json_object).toJson(QJsonDocument::Compact);

    QNetworkReply *reply = m_networkManager.post(request, post_data);
    connect(reply, &QNetworkReply::finished,[reply, this] {
       reply->deleteLater();

        if (reply->error() != QNetworkReply::NoError) {
            qWarning()<< "Network Error: "<< reply->errorString();
            return;
        }

        int StatusCode = reply->attribute(QNetworkRequest::HttpStatusCodeAttribute).toInt();

        if (StatusCode >= 200 && StatusCode < 300) {

            const QByteArray response_data = reply->readAll();

            QJsonParseError parse_error;

            QJsonDocument json_document = QJsonDocument::fromJson(response_data,&parse_error);

            if (parse_error.error == QJsonParseError::NoError && json_document.isObject()) {

                const QJsonObject root_object = json_document.object();

                token = root_object.value("token").toString();
                qDebug()<<"Token: "<<token;
            }
            else {
                qWarning()<<"Failed to parse json: "<<parse_error.errorString();
            }
        }
        else {
            qWarning()<<"Http Error: "<<reply->errorString();
        }
    });
}

void ConnectionHandler::userConnectionWebsocket() {

}

QNetworkRequest ConnectionHandler::createAuthRequest(const QString &endpoint) {
    QNetworkRequest request(m_baseURL);
    const QByteArray token;

    request.setRawHeader("Authorization", "Bearer" + token);

}

QNetworkRequest ConnectionHandler::createUnauthRequest(const QString &endpoint) {
    QNetworkRequest request(m_baseURL+endpoint);
    request.setRawHeader("Content-Type", "application/json");
    request.setTransferTimeout(5000);

    return request;
}
