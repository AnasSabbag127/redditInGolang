package models

import "github.com/google/uuid"

type Post struct {
	Id        uuid.UUID `json:"id"`
	PostTitle string    `json:"post_title"`
	PostText  string    `json:"post_text"`
	UserId    uuid.UUID `json:"user_id"`
}

type InputPost struct {
	PostTitle string    `json:"post_title"`
	PostText  string    `json:"post_text"`
	UserId    uuid.UUID `json:"user_id"`
}
