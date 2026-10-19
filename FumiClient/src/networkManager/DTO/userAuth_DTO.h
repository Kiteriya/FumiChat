//
// Created by Yukiteriya on 16.10.2026.
//

#ifndef FUMICLIENT_USERAUTH_DTO_H
#define FUMICLIENT_USERAUTH_DTO_H

#include <QString>

struct UserRegistration{
    QString username;
    QString password;
};

struct UserAuthorization {
    QString username;
    QString password;

};


#endif //FUMICLIENT_USERAUTH_DTO_H
