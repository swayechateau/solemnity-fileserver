package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/", DocumentHandler).Methods("GET")
	mux.HandleFunc("/api/v1/demo", DemoUploadHandler).Methods("GET")

	// mux.HandleFunc("/api/v1/upload", UploadHandler).Methods("POST")
	// mux.HandleFunc("/api/v1/view/{slug}", ViewFileHandler).Methods("GET")

	fmt.Println("Starting front end service on port 80")
	log.Panic(http.ListenAndServe(":80", mux))
}
