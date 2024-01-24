package models

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type FileAccess struct {
	Id           [16]byte  `json:"-" db:"Id"`
	FileId       [16]byte  `json:"-" db:"FileId"`
	Organization string    `db:"Organization"`
	Owner        string    `db:"AccessOwner"`
	IsPublic     bool      `db:"IsPublic"`
	Slug         string    `db:"Uri"`
	ShareCode    string    `db:"ShareCode"`
	AccessCode   string    `db:"AccessCode"`
	CreatedAt    time.Time `db:"CreatedAt"`
	UpdatedAt    time.Time `db:"UpdatedAt"`
}

type PublicFileAccess struct {
	FileId        [16]byte `json:"-" db:"FileId"`
	FileName      string   `db:"FileName"`
	FileExtension string   `db:"FileExtension"`
	FileType      string   `db:"FileType"`
	FileSize      int64    `db:"FileSize"`
	Uri           string   `db:"Uri"`
}

func (fa *FileAccess) GenerateSlug() {
	fa.Slug = uuid.New().String()
}

func (fa *FileAccess) GenerateShareCode() {
	fa.ShareCode = "s" + fa.Owner + generateCode(10)
}

func (fa *FileAccess) GenerateAccessCode() {
	fa.AccessCode = "a" + fa.Owner + generateCode(10)
}

func (fa *FileAccess) Get(c *pgx.Conn) error {
	ctx := context.Background()
	return c.QueryRow(ctx, "SELECT Id, FileId, Organization, AccessOwner, IsPublic, Uri, ShareCode, AccessCode, CreatedAt, UpdatedAt FROM FileAccess WHERE Uri = $1", fa.Slug).Scan(&fa.Id, &fa.FileId, &fa.Organization, &fa.Owner, &fa.IsPublic, &fa.Slug, &fa.ShareCode, &fa.AccessCode, &fa.CreatedAt, &fa.UpdatedAt)
}

func (fa *FileAccess) Create(c *pgx.Conn) error {
	ctx := context.Background()
	_, err := c.Exec(ctx, "INSERT INTO FileAccess (Id, FileId, AccessOwner, Organization, IsPublic, Uri, ShareCode, AccessCode) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", fa.Id, fa.FileId, fa.Owner, fa.Organization, fa.IsPublic, fa.Slug, fa.ShareCode, fa.AccessCode)
	return err
}

func (fa *FileAccess) Update(c *pgx.Conn) error {
	ctx := context.Background()
	_, err := c.Exec(ctx, "UPDATE FileAccess SET AccessOwner = $1, Organization = $2, IsPublic = $3, Uri = $4, ShareCode = $5, AccessCode = $6 WHERE Id=$7", fa.Owner, fa.Organization, fa.IsPublic, fa.Slug, fa.ShareCode, fa.AccessCode, fa.Id)
	return err
}

func (fa *FileAccess) Delete(c *pgx.Conn) error {
	ctx := context.Background()
	_, err := c.Exec(ctx, "DELETE FROM FileAccess WHERE Slug = $1", fa.Slug)
	return err
}

func (fa *FileAccess) GenerateId() {
	fa.Id = uuid.New()
}

func generateCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetPublicAccessFile(c *pgx.Conn, organization string) []PublicFileAccess {
	accessList := []PublicFileAccess{}
	ctx := context.Background()
	rows, err := c.Query(ctx, "SELECT FileId, Uri FROM FileAccess WHERE IsPublic = true AND Organization = $1", organization)
	if err != nil {
		log.Println(err)
		return accessList
	}
	defer rows.Close()

	for rows.Next() {
		var access PublicFileAccess
		err := rows.Scan(&access.FileId, &access.Uri)
		if err != nil {
			log.Println(err)
			return accessList
		}
		accessList = append(accessList, access)
	}
	// get file info
	for i, access := range accessList {
		file, err := GetFileById(c, access.FileId)
		if err != nil {
			log.Println(err)
			return accessList
		}
		accessList[i].FileName = file.Name
		accessList[i].FileExtension = file.Extension
		accessList[i].FileType = file.MimeType
		accessList[i].FileSize = file.Size
	}

	return accessList
}
