package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", rootHandler)
}

// display the api documentation
func rootHandler(w http.ResponseWriter, r *http.Request) {

}

// view public files
func publicHandler(w http.ResponseWriter, r *http.Request) {
	
}

// upload file/s
func uploadHandler(w http.ResponseWriter, r *http.Request) {

}

// download file/s
func downloadHandler(w http.ResponseWriter, r *http.Request) {

}

// delete file/s
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	
}	

// share file/s
func shareHandler(w http.ResponseWriter, r *http.Request) {
	
}

// access file/s
func accessHandler(w http.ResponseWriter, r *http.Request) {

}

// recover file/s
func recoveryHandler(w http.ResponseWriter, r *http.Request) {
	
}



struct File {
        'name',
        'path',
        'full_path',
        'size',
        'type',
        'extension'

}

struct Access {

        'organisation','owner','type',
        'public','slug','share_code',
        'access_code','file_id'
}

struct Recovery {
	'email', 
	'ip', 
	'domain',
	'code', 
	'access_id'
}

Schema::create('accesses', function (Blueprint $table) {
	$table->increments('id');
	$table->string('organisation')->default('global');
	$table->string('owner')->nullable();
	$table->string('type')->nullable();
	$table->boolean('public')->default(true);
	$table->string('slug')->unique();
	$table->string('share_code')->unique();
	$table->string('access_code')->unique();
	$table->unsignedInteger('file_id');
	$table->softDeletes();
	$table->timestamps();
});

Schema::create('files', function (Blueprint $table) {
	$table->increments('id');
	$table->string('name');
	$table->string('path');
	$table->string('full_path');
	$table->integer('size');
	$table->string('type');
	$table->string('extension');
	$table->softDeletes();
	$table->timestamps();
});

