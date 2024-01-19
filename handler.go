package main

import (
	"fmt"
	"io"
	"net/http"
)

// display the api documentation
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, documentation)
}

// display the api documentation
func demoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, uploadForm)
}

// upload file/s
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	owner := r.FormValue("owner")
	files := r.MultipartForm.File["files"]
	// Buffer for MIME type detection
	buffer := make([]byte, 512)

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hash, err := hashFileBlake2(file)
		if err != nil {
			file.Close()
			http.Error(w, "Failed to hash file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		access := FileAccess{FileHash: fmt.Sprintf("%x", hash), Owner: owner}

		// Check if file already exists
		if _, exists := savedFileHashes[access.FileHash]; exists {
			file.Close()
			fmt.Fprintf(w, "File %s skipped, already exists.\n", fileHeader.Filename)
			continue
		}

		// Save hash to the map
		savedFileHashes[access.FileHash] = true

		// Reset the file pointer before saving
		file.Seek(0, io.SeekStart)

		processedFile, err := convertFile(fileHeader, file, access.FileHash, buffer)
		if err != nil {
			file.Close()
			http.Error(w, "Failed to process file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()
		processedFile.Hash = access.FileHash

		fmt.Fprintf(w, "File %s uploaded successfully by %s with MIME type=%s and size=%d bytes. Hash: %s\n", processedFile.Name, access.Owner, processedFile.MimeType, processedFile.Size, processedFile.Hash)
	}

	fmt.Fprintf(w, "Upload successful")
}

// view public files
func publicHandler(w http.ResponseWriter, r *http.Request) {

}

// // download file/s
// func downloadHandler(w http.ResponseWriter, r *http.Request) {

// }

// // delete file/s
// func deleteHandler(w http.ResponseWriter, r *http.Request) {

// }

// // share file/s
// func shareHandler(w http.ResponseWriter, r *http.Request) {

// }

// // access file/s
// func accessHandler(w http.ResponseWriter, r *http.Request) {

// }

// // recover file/s
// func recoveryHandler(w http.ResponseWriter, r *http.Request) {

// }
