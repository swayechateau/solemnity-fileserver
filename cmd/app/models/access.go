package models

import "time"

type Access struct {
	Id int `json:"id"`
	// organisation is the organisation or service that owns the file
	Organisation string `json:"organisation"`
	// owner is the user who uploaded the file
	Owner      string    `json:"owner"`
	Type       string    `json:"type"`
	Public     bool      `json:"public"`
	Slug       string    `json:"slug"`
	ShareCode  string    `json:"share_code"`
	AccessCode string    `json:"access_code"`
	FileId     string    `json:"file_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
