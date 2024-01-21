package models

import "time"

type File struct {
	// file uuid
	Id string `json:"id"`
	// file hash
	Hash      string    `json:"hash"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	FullPath  string    `json:"full_path"`
	Size      string    `json:"size"`
	Type      string    `json:"type"`
	Extension string    `json:"extension"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrganisationFiles struct {
	Organisation string    `json:"organisation_id"`
	FileId       string    `json:"file_id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserFiles struct {
	UserId    string    `json:"user_id"`
	FileId    string    `json:"file_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (file *File) GenerateHash() {
	// generate hash
}
