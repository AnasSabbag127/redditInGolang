package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
}

type InputUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
