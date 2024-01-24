package internal

import (
	"fileserver/internal/fcrypt"
	"fileserver/internal/models"
	"fileserver/internal/templates"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/swayedev/way"
)

// cipherKey is the key used only for development and testing
var cipherKey = []byte("63ff16e308e2f50e2b9be6565e37de9b1dbc81e7582471487704603ddf9a1288")
var maxMemory int64 = 1024 * 1024 * 10 // 10 MB

// display the api documentation
func RootHandler(c *way.Context) {
	c.HTML(200, templates.Documentation)
}

// display the api documentation
func DemoHandler(c *way.Context) {
	c.HTML(200, templates.UploadForm)
}
func UploadWithOptionalEncryptionHandler(c *way.Context) {
	switch c.Request.FormValue("encrypt") {
	case "true":
		UploadAndEncryptHandler(c)
	default:
		UploadHandler(c)
	}
}

// upload file/s
func UploadHandler(c *way.Context) {
	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		log.Printf("Error:%s", err)
		http.Error(c.Response, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Are there any files?
	files := c.Request.MultipartForm.File["media"]
	if len(files) == 0 {
		http.Error(c.Response, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Get the owner, isPublic, and organization from the form
	owner := c.Request.FormValue("owner")
	if owner == "" {
		// was anonymous
		owner = "public" // default to public
	}

	organization := c.Request.FormValue("organization")
	if organization == "" {
		organization = "global"
	}

	isPublic := true
	if c.Request.FormValue("public") != "" {
		isPublic, _ = strconv.ParseBool(c.Request.FormValue("public"))
	}
	accessList := []models.FileAccess{}
	// Process each file
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error:%s", err)
			http.Error(c.Response, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hash, err := fcrypt.HashWithBlake2(file)
		if err != nil {
			log.Printf("Error:%s", err)
			file.Close()
			http.Error(c.Response, "Failed to hash file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fileHash := fmt.Sprintf("%x", hash)

		// Check if the files already exist
		fileId, err := models.GetFileId(c.GetDB().Pgx(), organization, fileHash)

		if err != nil && err != pgx.ErrNoRows {
			log.Printf("Error:%s", err)
			file.Close()
			http.Error(c.Response, "Failed to get file id: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err == pgx.ErrNoRows {
			// Reset the file pointer before saving
			file.Seek(0, io.SeekStart)
			f, err := models.StoreFile(fileHeader, file, fileHash, organization)
			if err != nil {
				log.Printf("Error:%s", err)
				file.Close()
				http.Error(c.Response, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			file.Close()
			fileId = f.Id
			// save file to db
			_, err = c.PgxExec(c.Request.Context(), "INSERT INTO Files (Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", f.Id, f.Organization, f.Name, f.Extension, f.MimeType, f.Size, f.Hash, f.Path, f.FullPath)
			if err != nil {
				log.Printf("Error:%s", err)
				http.Error(c.Response, "Failed to save file to database: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		access := models.FileAccess{
			FileId:       fileId,
			Organization: organization,
			Owner:        owner,
			IsPublic:     isPublic,
		}

		access.GenerateId()
		access.GenerateSlug()
		access.GenerateShareCode()
		access.GenerateAccessCode()

		if err := access.Create(c.GetDB().Pgx()); err != nil {
			log.Printf("Error:%s", err)
			http.Error(c.Response, "Failed to save file access to database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		accessList = append(accessList, access)
	}

	c.JSON(200, accessList)
}

func UploadAndEncryptHandler(c *way.Context) {
	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		http.Error(c.Response, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Are there any files?
	files := c.Request.MultipartForm.File["media"]
	if len(files) == 0 {
		http.Error(c.Response, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Get the owner, isPublic, and organization from the form
	owner := c.Request.FormValue("owner")
	if owner == "" {
		// was anonymous
		owner = "public" // default to public
	}

	organization := c.Request.FormValue("organization")
	if organization == "" {
		organization = "global"
	}

	isPublic := true
	if c.Request.FormValue("public") != "" {
		isPublic, _ = strconv.ParseBool(c.Request.FormValue("public"))
	}

	// Process each file
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(c.Response, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hash, err := fcrypt.HashWithBlake2(file)
		if err != nil {
			file.Close()
			http.Error(c.Response, "Failed to hash file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fileHash := fmt.Sprintf("%x", hash)

		// Check if the files already exist
		fileId, err := models.GetFileId(c.GetDB().Pgx(), organization, fileHash)

		if err != nil && err != pgx.ErrNoRows {
			file.Close()
			http.Error(c.Response, "Failed to get file id: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err == pgx.ErrNoRows {
			// Reset the file pointer before saving
			file.Seek(0, io.SeekStart)
			f, err := models.StoreAndEncryptFile(fileHeader, file, fileHash, organization, cipherKey)
			if err != nil {
				file.Close()
				http.Error(c.Response, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			file.Close()
			fileId = f.Id
			// save file to db
			_, err = c.PgxExec(c.Request.Context(), "INSERT INTO Files (Id, Organization, FileName, FileExtension, FileType, FileSize, FileHash, FilePath, FileFullPath) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", f.Id, f.Organization, f.Name, f.Extension, f.MimeType, f.Size, f.Hash, f.Path, f.FullPath)
			if err != nil {
				http.Error(c.Response, "Failed to save file to database: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		access := models.FileAccess{
			FileId:       fileId,
			Organization: organization,
			Owner:        owner,
			IsPublic:     isPublic,
		}

		access.GenerateId()
		access.GenerateSlug()
		access.GenerateShareCode()
		access.GenerateAccessCode()

		if access.Create(c.GetDB().Pgx()) != nil {
			http.Error(c.Response, "Failed to save file access to database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(c.Response, "File %s uploaded successfully by %s with MIME type=%s and size=%d bytes. Hash: %s\n", fileHeader.Filename, access.Owner, fileHeader.Header.Get("Content-Type"), fileHeader.Size, fileHash)
	}

	fmt.Fprintf(c.Response, "Upload successful")

}

// view or download file
func ViewHandler(c *way.Context) {
	// find access by slug
	slug := c.Parms()["slug"]
	access := models.FileAccess{Slug: slug}
	if err := access.Get(c.GetDB().Pgx()); err != nil {
		http.Error(c.Response, "Failed to get file access: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// get file by id
	f := models.File{Id: access.FileId}
	if err := f.Get(c.GetDB().Pgx()); err != nil {
		http.Error(c.Response, "Failed to get file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// if access is public, serve file
	if access.IsPublic {
		// Stream the file
		f.Serve(c.Response, c.Request, cipherKey)
		return
	}

	// if access code or share code is valid, serve file
	if c.Request.FormValue("access_code") == access.AccessCode || c.Request.FormValue("share_code") == access.ShareCode {
		// Stream the file
		f.Serve(c.Response, c.Request, cipherKey)
		return
	}
	// if access code or share code is invalid, return 403
	http.Error(c.Response, "Access denied", http.StatusForbidden)
}

// view public files
func PublicHandler(c *way.Context) {
	// get all public access records
	organization := c.Request.FormValue("organization")
	if organization == "" {
		organization = "global"
	}
	c.JSON(200, models.GetPublicAccessFile(c.GetDB().Pgx(), organization))
}
