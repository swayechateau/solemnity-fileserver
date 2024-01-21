package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swayedev/way"
)

// display the api documentation
func rootHandler(c *way.Context) {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(c.Response, documentation)
}

// display the api documentation
func demoHandler(c *way.Context) {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(c.Response, uploadForm)
}

// upload file/s
func uploadHandler(c *way.Context) {
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

		hash, err := hashFileBlake2(file)
		if err != nil {
			file.Close()
			http.Error(c.Response, "Failed to hash file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		access := FileAccess{FileHash: fmt.Sprintf("%x", hash), Owner: owner}

		// Check if file already exists
		if _, exists := savedFileHashes[access.FileHash]; exists {
			file.Close()
			fmt.Fprintf(c.Response, "File %s skipped, already exists.\n", fileHeader.Filename)
			continue
		}

		// Save hash to the map
		savedFileHashes[access.FileHash] = true

		// Reset the file pointer before saving
		file.Seek(0, io.SeekStart)

		processedFile, err := convertFile(fileHeader, file, access.FileHash, buffer)
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

// view public files
func publicHandler(c *way.Context) {
	// search for public files under the global organisation
	// return Access::where('public', 1)->get();
}

// // download file/s
func viewFileHandler(c *way.Context) {
	slug := mux.Vars(c.Request)["slug"]
	fmt.Fprintf(c.Response, "Viewing file %s", slug)
}

// // delete file/s
// func deleteHandler(c *way.Context) {

// }

// // share file/s
// func shareHandler(c *way.Context) {

// }

// // access file/s
// func accessHandler(c *way.Context) {

// }

// // recover file/s
// func recoveryHandler(c *way.Context) {

// }

func uploadAndEncryptHandler(c *way.Context) {
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

		hash, err := hashFileBlake2(file)
		if err != nil {
			file.Close()
			http.Error(c.Response, "Failed to hash file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		access := FileAccess{FileHash: fmt.Sprintf("%x", hash), Owner: owner}

		// Check if file already exists
		if _, exists := savedFileHashes[access.FileHash]; exists {
			file.Close()
			fmt.Fprintf(c.Response, "File %s skipped, already exists.\n", fileHeader.Filename)
			continue
		}

		// Save hash to the map
		savedFileHashes[access.FileHash] = true

		// Reset the file pointer before saving
		file.Seek(0, io.SeekStart)

		processedFile, err := convertFile(fileHeader, file, access.FileHash, buffer)
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
