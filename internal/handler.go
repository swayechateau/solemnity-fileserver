package internal

import (
	"fileserver/internal/fcrypt"
	"fileserver/internal/models"
	"fileserver/internal/templates"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/swayedev/way"
)

var cipherKey = []byte("63d2c76bcaad9fcb757f57458d7c6384")
var maxMemory int64 = 1024 * 1024 * 10 // 10 MB

// display the api documentation
func RootHandler(c *way.Context) {
	c.HTML(200, templates.Documentation)
}

// display the api documentation
func DemoHandler(c *way.Context) {
	c.HTML(200, templates.UploadForm)
}

// upload file/s
func UploadHandler(c *way.Context) {
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
		fileId, err := models.GetFileId(c, organization, fileHash)

		if err != nil && err != pgx.ErrNoRows {
			file.Close()
			http.Error(c.Response, "Failed to get file id: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err == pgx.ErrNoRows {
			// Reset the file pointer before saving
			file.Seek(0, io.SeekStart)
			f, err := models.StoreFile(fileHeader, file, fileHash, organization)
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

		if access.Create(c) != nil {
			http.Error(c.Response, "Failed to save file access to database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// // Check if file already exists
		// if exists, _ := CheckFile(c, access.FileHash); exists {
		// 	file.Close()

		// 	fmt.Fprintf(c.Response, "File %s skipped, already exists.\n", fileHeader.Filename)
		// 	continue
		// }

		// // Save hash to the map
		// savedFileHashes[access.FileHash] = true

		// Save the file to the database

		// fmt.Fprintf(c.Response, "File %s uploaded successfully by %s with MIME type=%s and size=%d bytes. Hash: %s\n", f.Name, access.Owner, f.MimeType, f.Size, f.Hash)
		fmt.Fprintf(c.Response, "File %s uploaded successfully by %s with MIME type=%s and size=%d bytes. Hash: %s\n", fileHeader.Filename, access.Owner, fileHeader.Header.Get("Content-Type"), fileHeader.Size, fileHash)
	}

	fmt.Fprintf(c.Response, "Upload successful")
}

// View or download file/s
func ViewFileHandler(c *way.Context) {
	// find access by slug
	slug := c.Parms()["slug"]
	access := models.FileAccess{Slug: slug}
	if err := access.Get(c); err != nil {
		http.Error(c.Response, "Failed to get file access: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// get file by id
	f := models.File{Id: access.FileId}
	if err := f.Get(c); err != nil {
		http.Error(c.Response, "Failed to get file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// if access is public, serve file
	if access.IsPublic {
		// Stream the file
		http.ServeFile(c.Response, c.Request, f.FullPath)
		return
	}

	// if access code or share code is valid, serve file
	if c.Request.FormValue("access_code") == access.AccessCode || c.Request.FormValue("share_code") == access.ShareCode {
		// Stream the file
		http.ServeFile(c.Response, c.Request, f.FullPath)
		return
	}
	// if access code or share code is invalid, return 403
	http.Error(c.Response, "Access denied", http.StatusForbidden)
}

func UploadAndEncryptHandler(c *way.Context) {
	if err := c.Request.ParseMultipartForm(maxMemory); err != nil {
		http.Error(c.Response, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	owner := c.Request.FormValue("owner")
	files := c.Request.MultipartForm.File["files"]
	// Buffer for MIME type detection
	buffer := make([]byte, 512)

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
		access := models.FileAccess{FileHash: fmt.Sprintf("%x", hash), Owner: owner}

		// check db for file hash
		// if exists, skip
		// if !exists, continue

		// // Check if file already exists
		// if _, exists := savedFileHashes[access.FileHash]; exists {
		// 	file.Close()
		// 	fmt.Fprintf(c.Response, "File %s skipped, already exists.\n", fileHeader.Filename)
		// 	continue
		// }

		// // Save hash to the map
		// access
		// savedFileHashes[access.FileHash] = true

		// Reset the file pointer before saving
		file.Seek(0, io.SeekStart)

		processedFile, err := models.StoreAndEncryptFile(fileHeader, file, access.FileHash, buffer, cipherKey)
		if err != nil {
			file.Close()
			http.Error(c.Response, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()
		processedFile.Hash = access.FileHash

		fmt.Fprintf(c.Response, "File %s uploaded successfully by %s with MIME type=%s and size=%d bytes. Hash: %s\n", processedFile.Name, access.Owner, processedFile.MimeType, processedFile.Size, processedFile.Hash)
	}

	fmt.Fprintf(c.Response, "Upload successful")
}

func ViewAndDecryptFileHandler(c *way.Context) {
	// Extract file hash from the URL
	slug := c.Parms()["slug"]
	fileHash := slug

	// Get the path to the encrypted file using the fileHash
	// For this example, we'll assume the file path is "./uploads/" + fileHash
	// In a real application, you might look up a database or a map
	encryptedFilePath := "./uploads/" + fileHash
	decryptedFilePath := "./temp/decrypted_" + fileHash // Temporary path for decrypted file

	// Define your AES key
	key := cipherKey // Replace with your key

	// Decrypt the file
	if err := fcrypt.DecryptWithGCM(encryptedFilePath, decryptedFilePath, key); err != nil {
		http.Error(c.Response, "Failed to decrypt file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Stream the decrypted file
	http.ServeFile(c.Response, c.Request, decryptedFilePath)

	// Optionally delete the temporary decrypted file after serving
	os.Remove(decryptedFilePath)
}
