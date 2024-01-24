package models

import (
	"context"
	"fileserver/internal/fcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type File struct {
	Id           [16]byte  `db:"Id"`
	Organization string    `db:"Organization"`
	IsEncrypted  bool      `db:"IsEncrypted"`
	Name         string    `db:"FileName"`
	Extension    string    `db:"FileExtension"`
	MimeType     string    `db:"FileType"`
	Size         int64     `db:"FileSize"`
	Hash         string    `db:"FileHash"`
	Path         string    `db:"FilePath"`
	FullPath     string    `db:"FileFullPath"`
	CreatedAt    time.Time `db:"CreatedAt"`
	UpdatedAt    time.Time `db:"UpdatedAt"`
}

func (f *File) Create(fileHeader *multipart.FileHeader, file multipart.File, hash string, organization string) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}
	// dir, _ := os.Getwd()
	f.Id = uuid.New()
	f.Name = fileHeader.Filename
	f.IsEncrypted = false
	f.Organization = checkOrganization(organization)
	f.Extension = filepath.Ext(fileHeader.Filename)
	f.MimeType = http.DetectContentType(buffer)
	f.Hash = hash
	f.Path = "./uploads/" + f.Organization + "/"
	// f.Path = "/" + dir + "/uploads/" + f.Organization + "/"
	f.FullPath = f.Path + f.Hash + f.Extension
	// Reset file read pointer
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	f.Size = fileHeader.Size

	dst, err := os.Create(f.FullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func (f *File) CreateAndEncrypt(fileHeader *multipart.FileHeader, file multipart.File, hash string, organization string, key []byte) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}
	dir, _ := os.Getwd()
	f.Id = uuid.New()
	f.Name = fileHeader.Filename
	f.IsEncrypted = true
	f.Organization = checkOrganization(organization)
	f.Extension = filepath.Ext(fileHeader.Filename)
	f.MimeType = http.DetectContentType(buffer)
	f.Hash = hash
	f.Path = dir + "/uploads/" + f.Organization + "/"
	f.FullPath = "/" + f.Path + f.Hash + f.Extension
	// Reset file read pointer
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	f.Size = fileHeader.Size

	dst, err := os.Create(f.FullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if err = fcrypt.EncryptWithGCM(file, dst, key); err != nil {
		return err
	}

	return nil
}

func (f *File) Serve(w http.ResponseWriter, r *http.Request, key []byte) {
	if f.IsEncrypted {
		f.DecryptAndServe(w, r, key)
		return
	}
	http.ServeFile(w, r, f.FullPath)
}

func (f *File) DecryptAndServe(w http.ResponseWriter, r *http.Request, key []byte) {
	// Temporary path for decrypted file
	decryptedFilePath := "./uploads/temp/decrypted_" + f.Hash + f.Extension
	// Decrypt the file
	if err := fcrypt.DecryptWithGCM(f.FullPath, decryptedFilePath, key); err != nil {
		http.Error(w, "Failed to decrypt file: "+err.Error(), http.StatusInternalServerError)
		os.Remove(decryptedFilePath)
		return
	}
	// Stream the decrypted file
	http.ServeFile(w, r, decryptedFilePath)
	// Optionally delete the temporary decrypted file after serving
	os.Remove(decryptedFilePath)
}

func StoreFile(fileHeader *multipart.FileHeader, file multipart.File, hash string, organization string) (File, error) {
	f := File{}
	if err := f.Create(fileHeader, file, hash, organization); err != nil {
		return f, err
	}
	return f, nil
}

func StoreAndEncryptFile(fileHeader *multipart.FileHeader, file multipart.File, hash string, organization string, key []byte) (File, error) {
	f := File{}
	if err := f.CreateAndEncrypt(fileHeader, file, hash, organization, key); err != nil {
		return f, err
	}
	return f, nil
}

func (f *File) Get(c *pgx.Conn) error {
	ctx := context.Background()
	return c.QueryRow(ctx, "SELECT Id, Organization, IsEncrypted, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Id=$1", f.Id).Scan(&f.Id, &f.Organization, &f.IsEncrypted, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
}

func (f *File) Save(c *pgx.Conn) error {
	ctx := context.Background()
	_, err := c.Exec(ctx, "INSERT INTO Files (Id, Organization, IsEncrypted, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", f.Id, f.Organization, f.IsEncrypted, f.Name, f.Extension, f.MimeType, f.Size, f.Hash, f.Path, f.FullPath)
	return err
}

func (f *File) Update(c *pgx.Conn) error {
	ctx := context.Background()
	_, err := c.Exec(ctx, "UPDATE Files SET Organization = $1, FileName = $2, FileExtension = $3, FileType = $4, FileSize = $5, FileHash = $6, FilePath = $7, FileFullPath = $8 IsEncrypted =$9 WHERE Id=$10", f.Organization, f.Name, f.Extension, f.MimeType, f.Size, f.Hash, f.Path, f.FullPath, f.IsEncrypted, f.Id)
	return err
}

func GetFileId(c *pgx.Conn, organization string, hash string) ([16]byte, error) {
	var id [16]byte
	ctx := context.Background()
	row := c.QueryRow(ctx, `SELECT Id FROM Files WHERE Organization=$1 AND FileHash=$2;`, checkOrganization(organization), hash)
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func GetFileById(c *pgx.Conn, id [16]byte) (File, error) {
	f := File{}
	ctx := context.Background()
	err := c.QueryRow(ctx, `SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Id=$1;`, id).Scan(&f.Id, &f.Organization, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return f, err
	}

	return f, nil
}

func GetFile(c *pgx.Conn, organization string, hash string) (File, error) {
	f := File{}
	ctx := context.Background()
	err := c.QueryRow(ctx, `SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Organization=$1 AND FileHash=$2;`, checkOrganization(organization), hash).Scan(&f.Id, &f.Organization, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return f, err
	}

	return f, nil
}

func GetFiles(c *pgx.Conn, organization string) ([]File, error) {
	var files []File
	ctx := context.Background()
	rows, err := c.Query(ctx, `SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Organization=$1;`, checkOrganization(organization))
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		f := File{}
		err := rows.Scan(&f.Id, &f.Organization, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return files, err
		}
		files = append(files, f)
	}

	return files, nil
}

func checkOrganization(org string) string {
	if org == "" {
		return "global"
	}
	return org
}
