package models

import (
	"fileserver/internal/fcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/swayedev/way"
)

type File struct {
	Id           [16]byte
	Organization string
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

func (f *File) Get(w *way.Context) error {
	ctx := w.Request.Context()
	return w.PgxQueryRow(ctx, "SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Id=$1", f.Id).Scan(&f.Id, &f.Organization, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
}

func (f *File) Create(fileHeader *multipart.FileHeader, file multipart.File, hash string, organization string) error {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}
	dir, _ := os.Getwd()
	f.Id = uuid.New()
	f.Name = fileHeader.Filename
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

	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func (f *File) CreateAndEncrypt(fileHeader *multipart.FileHeader, file multipart.File, hash string, buffer []byte, key []byte) error {
	_, err := file.Read(buffer)
	if err != nil {
		return err
	}
	dir, _ := os.Getwd()
	f.Name = fileHeader.Filename
	f.Extension = filepath.Ext(fileHeader.Filename)
	f.MimeType = http.DetectContentType(buffer)
	f.Hash = hash
	f.Path = dir + "/uploads"
	f.FullPath = "/" + f.Path + f.Hash + f.Extension
	// Reset file read pointer
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	f.Size = fileHeader.Size

	dst, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if err = fcrypt.EncryptWithGCM(file, dst, key); err != nil {
		return err
	}

	return nil
}

func (f *File) Save(w *way.Context) error {
	ctx := w.Request.Context()
	return w.PgxExecNoResult(ctx, "INSERT INTO Files (Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", f.Id, f.Organization, f.Name, f.Extension, f.MimeType, f.Size, f.Hash, f.Path, f.FullPath)
}

func (f *File) Update(w *way.Context) error {
	ctx := w.Request.Context()
	return w.PgxExecNoResult(ctx, "UPDATE Files SET Organization = $1, FileName = $2, FileExtension = $3, FileType = $4, FileSize = $5, FileHash = $6, FilePath = $7, FileFullPath = $8 WHERE Id=$9", f.Organization, f.Name, f.Extension, f.MimeType, f.Size, f.Hash, f.Path, f.FullPath, f.Id)
}

func StoreFile(fileHeader *multipart.FileHeader, file multipart.File, hash string, organization string) (File, error) {
	f := File{}
	if err := f.Create(fileHeader, file, hash, organization); err != nil {
		return f, err
	}
	return f, nil
}

func StoreAndEncryptFile(fileHeader *multipart.FileHeader, file multipart.File, hash string, buffer []byte, key []byte) (File, error) {
	f := File{}
	if err := f.CreateAndEncrypt(fileHeader, file, hash, buffer, key); err != nil {
		return f, err
	}
	return f, nil
}

func GetFileId(w *way.Context, organization string, hash string) ([16]byte, error) {
	var id [16]byte

	err := w.PgxQueryRow(w.Request.Context(), `SELECT Id FROM Files WHERE Organization=$1 AND FileHash=$2;`, checkOrganization(organization), hash).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func GetFileById(w *way.Context, id [16]byte) (File, error) {
	f := File{}

	err := w.PgxQueryRow(w.Request.Context(), `SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE FileID=$2;`, id).Scan(&f.Id, &f.Organization, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return f, err
	}

	return f, nil
}

func GetFile(w *way.Context, organization string, hash string) (File, error) {
	f := File{}

	err := w.PgxQueryRow(w.Request.Context(), `SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Organization=$1 AND FileHash=$2;`, checkOrganization(organization), hash).Scan(&f.Id, &f.Organization, &f.Name, &f.Extension, &f.MimeType, &f.Size, &f.Hash, &f.Path, &f.FullPath, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return f, err
	}

	return f, nil
}

func GetFiles(w *way.Context, organization string) ([]File, error) {
	var files []File

	rows, err := w.PgxQuery(w.Request.Context(), `SELECT Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath, CreatedAt, UpdatedAt FROM Files WHERE Organization=$1;`, checkOrganization(organization))
	if err != nil {
		return files, err
	}

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
