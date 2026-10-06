package models

import "uuid"

type User struct {
	Id       uuid.UUID
	Username string
	Password string
}
