package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
)

type File struct {
	Name      string
	Extension string
	MimeType  string
	Size      int64
	Hash      string
	Path      string
	FullPath  string
}

type FileAccess struct {
	FileHash   string
	Owner      string
	Provider   string
	IsPublic   bool
	Slug       string
	ShareCode  string
	AccessCode string
}

type Recovery struct {
	Email    string
	IP       string
	Domain   string
	Code     string
	AccessId string
}

func hashFileBlake2(file multipart.File) ([]byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func convertFile(fileHeader *multipart.FileHeader, file multipart.File, hash string, buffer []byte) (File, error) {
	f := File{}

	_, err := file.Read(buffer)
	if err != nil {
		return f, err
	}
	f.Name = fileHeader.Filename
	f.Extension = filepath.Ext(fileHeader.Filename)
	f.MimeType = http.DetectContentType(buffer)
	f.Hash = hash

	// Reset file read pointer
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return f, err
	}

	f.Size = fileHeader.Size
	f.Path = "./uploads/" + f.Hash + f.Extension
	f.FullPath, _ = os.Getwd()
	f.FullPath += "/" + f.Path
	dst, err := os.Create(f.Path)
	if err != nil {
		return f, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return f, err
	}

	return f, nil
}
