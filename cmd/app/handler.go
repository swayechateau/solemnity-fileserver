package main

import "net/http"

// view API Doc
func DocumentHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/index.html")
}

func DemoUploadHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/upload.html")
}

// // Upload File
// func UploadHandler(w http.ResponseWriter, r *http.Request) {
// 	// upload file
// 	organisation := r.FormValue("organisation");
// 	owner := r.FormValue("owner");
// 	Ftype := r.FormValue("type");
// 	public := r.FormValue("public");
// 	date := r.FormValue("date");
// 	if(r.FormValue("multi") == "true") {
// 		multiFileUpload(request->file('media'), $organisation, $owner, $type, $public, $date);
// 		return
// 	}

// 	upload($request->file('media'), $organisation, $owner, $type, $public, $date);

// }

// // view all public files
// func ViewPublicFilesHandler(w http.ResponseWriter, r *http.Request) {
// 	// view public files
// 	// return Access::where('public', 1)->get();
// }

// // View File
// func ViewFileHandler(w http.ResponseWriter, r *http.Request) {
// 	// view file
// 	$access = Access::where('slug', $slug)->first();
// 	if($access) {
// 		if(Helper::checkAccess($access, $request)) {
// 			$file = File::find($access->file_id);
// 			if (file_exists($file->full_path)) {
// 				header("Content-Type: ".$file->type);
// 				header("Content-Length:".$file->size);
// 				header('Content-Disposition: inline; filename="' . $file->name . '"');
// 				header('Content-Transfer-Encoding: binary');
// 				header('Accept-Ranges: bytes');
// 				$fp = fopen($file->full_path, 'rb');
// 				return fpassthru($fp);
// 			} else {
// 				return response()->json(['Error'=> 'File not found.'], 503);
// 			}
// 		}
// 		return response()->json(['error'=> 'Invalid Access Rights'],401);
// 	}
// 	return response()->json(['error' => 'Nothing found!'],404);
// }

// // View File Details
// func ViewFileDetailsHandler(w http.ResponseWriter, r *http.Request) {
// 	// view file details
// 	$access = Access::where('slug', $slug)->first();
// 	if($access) {
// 		if(Helper::checkAccessCode($access->access_code, $request)) {
// 			$file = File::find($access->file_id);
// 			if (file_exists($file->full_path)) {
// 				return [
// 					'access' => $access,
// 					'file' => $file,
// 				];;
// 			} else {
// 				return response()->json(['Error'=> 'File not found.'], 503);
// 			}
// 		}
// 		return response()->json(['error'=> 'Invalid Access Rights'],401);
// 	}
// 	return response()->json(['error' => 'Nothing found!'],404);
// }

// // Modify File
// func ModifyFileHandler(w http.ResponseWriter, r *http.Request) {
// 	// modify file
// 	$access = Access::where('slug', $slug)->first();
// 	if($access) {
// 		if(Helper::checkAccessCode($access->access_code, $request)) {
// 			if(Helper::updateAccess($slug, $request)) {
// 				return response()->json(['message' => 'File updated!'], 204);
// 			}
// 			return response()->json(['error' => 'An error Occured!'], 503);
// 		}
// 		return response()->json(['error'=> 'Invalid Access Rights'],401);
// 	}
// 	return response()->json(['error' => 'Nothing found!'],404);
// }

// // Remove File
// func RemoveFileHandler(w http.ResponseWriter, r *http.Request) {
// 	// remove file
// 	$access = Access::where('slug', $slug)->first();
// 	if($access) {
// 		if(Helper::checkAccessCode($access->access_code, $request)) {
// 			if(Helper::deleteAccess($slug)) {
// 				return response()->json(['message' => 'File deleted!'], 204);
// 			}
// 			return response()->json(['error' => 'An error Occured!'], 503);
// 		}
// 		return response()->json(['error'=> 'Invalid Access Rights'],401);
// 	}
// 	return response()->json(['error' => 'Nothing found!'],404);
// }

// func ViewUserFilesHandler(w http.ResponseWriter, r *http.Request) {
// 	// view user files
// }

// // Admin Only Routes
// func ViewAllFilesHandler(w http.ResponseWriter, r *http.Request) {
// 	// view all files
// 	// return Access::all();
// }

// func render(w http.ResponseWriter, t string) {

// 	var templateSlice []string
// 	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

// 	for _, x := range partials {
// 		templateSlice = append(templateSlice, x)
// 	}

// 	tmpl, err := template.ParseFiles(templateSlice...)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := tmpl.Execute(w, nil); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
