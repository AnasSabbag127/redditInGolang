package models

import "github.com/google/uuid"

type Comments struct {
	Id      uuid.UUID `json:"id"`
	PostId  uuid.UUID `json:"post_id"`
	UserId  uuid.UUID `json:"user_id"`
	Comment string    `json:"comment"`
}

type InputComment struct {
	PostId  uuid.UUID `json:"post_id"`
	UserId  uuid.UUID `json:"user_id"`
	Comment string    `json:"comment"`
}

type GetComment struct {
	UserId  uuid.UUID `json:"user_id"`
	Comment string    `json:"comment"`
}

type UpdateComment struct {
	Comment string `json:"comment"`
}
