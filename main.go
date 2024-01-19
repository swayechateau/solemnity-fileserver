package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var savedFileHashes = make(map[string]bool)
var maxMemory int64 = 1024 * 1024 * 10 // 10 MB

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/upload", demoHandler).Methods("GET")
	// API Routes
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/view/{fileHash}", viewFileHandler).Methods("GET")
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

// Schema::create('accesses', function (Blueprint $table) {
// 	$table->increments('id');
// 	$table->string('organisation')->default('global');
// 	$table->string('owner')->nullable();
// 	$table->string('type')->nullable();
// 	$table->boolean('public')->default(true);
// 	$table->string('slug')->unique();
// 	$table->string('share_code')->unique();
// 	$table->string('access_code')->unique();
// 	$table->unsignedInteger('file_id');
// 	$table->softDeletes();
// 	$table->timestamps();
// });

// Schema::create('files', function (Blueprint $table) {
// 	$table->increments('id');
// 	$table->string('name');
// 	$table->string('path');
// 	$table->string('full_path');
// 	$table->integer('size');
// 	$table->string('type');
// 	$table->string('extension');
// 	$table->softDeletes();
// 	$table->timestamps();
// });
