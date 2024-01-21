package main

// func uploadFile(file []byte, organisation string, owner string, fType string, public string, date string) {

//         if organisation == "" {
//             organisation = "global";
//         }

//         public = fileBoolean(public);

//         destinationPath := generateFilePath(getMimeType(), organisation, type, owner, date);
//         // check if file name exists
//         name = file->getClientOriginalName();
//         unique = false;
//         looped = 1;
//         if (file_exists(destinationPath . '/' . name)) {
//             while (unique === false) {
//                 if (file_exists(destinationPath . '/' . '(' . looped . ')_' . name)) {
//                     looped++;
//                 } else {
//                     unique = true;
//                 }
//             }
//             name = '(' . looped . ')_' . name;
//         }

//         // upload file
//         uploadedFile = File::create([
//             "name" => name,
//             "extension" => file->getClientOriginalExtension(),
//             "path" => destinationPath,
//             "full_path" => destinationPath . '/' . name,
//             "size" => file->getSize(),
//             "type" => file->getMimeType(),
//         ]);
//         file->move(destinationPath, name);

//         // create file access
//         slug = str_replace(' ', '-', organisation) . generateSlug() . uploadedFile->id;
//         shareCode = 's' . uploadedFile->id . generateCode();
//         accessCode = 'a' . uploadedFile->id . generateCode();
//         access = Access::create([
//             'organisation' => organisation,
//             'owner' => owner,
//             'type' => type,
//             'public' => public,
//             'slug' => slug,
//             'share_code' => shareCode,
//             'access_code' => accessCode,
//             'file_id' => uploadedFile->id,
//         ]);
//         // return file with access
//         if (access && uploadedFile) {
//             return [
//                 'access' => access,
//                 'file' => uploadedFile,
//             ];
//         }
//         return false;
//     }

//     func updateAccess(slug, request)
//     {
//         path = request->path;
//         name = request->name;
//         access = Access::where('slug', slug);
//         if (request->organisation) {
//             access->organisation = request->organisation;
//         }
//         if (request->owner) {
//             access->owner = request->owner;
//         }
//         if (request->type) {
//             access->type = request->type;
//         }
//         if (request->public) {
//             access->public = convertBoolean(request->public);
//         }

//         if (path || name) {
//             file = File::find(access->file_id);
//             updatedFile = updateFile(file, path, name);
//             if (updatedFile) {return true;}
//         }

//         return false;
//     }

//     func updateFile(file, cpath, cname)
//     {
//         if (cname || cpath) {
//             ext = file->extension;
//             path = file->path;
//             name = file->name;

//             // change path
//             if (cpath) {
//                 path = file->path;
//                 if (is_writable(path)) {
//                     name = checkFileName(path, name);
//                     file->path = path;
//                 }
//                 return false;
//             }

//             // check if name has .extention
//             // add .extention to name end if missing
//             if (cname) {
//                 name = cname;
//                 if (!strpos(name, ext)) {
//                     name = name . ext;
//                 }
//                 name = checkFileName(path, name);
//                 file->name = name;
//             }
//             // update file
//             rename(file->full_path, path . '/' . name);
//             // update full_path
//             file->full_path = path . '/' . name;
//             //update file record
//             return file->save();
//         }

//         return false;

//     }

//     func checkFileName(path, name)
//     {
//         unique = false;
//         looped = 1;
//         while (unique === false) {
//             if (file_exists(path . '/' . name)) {
//                 name = '(' . looped . ')_' . name;
//                 looped++;
//             } else {
//                 unique = true;
//             }
//         }
//         return name;
//     }

//     func checkAccess(access, request)
//     {

//         if (checkBoolean(access->public) === false) {
//             if (checkAccessCode(access->access_code, request)) {
//                 return true;
//             }

//             if (checkShareCode(access->share_code, request)) {
//                 return true;
//             }

//             return false;
//         }

//         return true;
//     }

//     func checkAccessCode(accessCode, request)
//     {

//         if (request->header('access_code')) {
//             return checkCode(accessCode, request->header('access_code'));
//         }

//         if (request->input('access_code')) {
//             return checkCode(accessCode, request->input('access_code'));
//         }

//         return false;

//     }

//     func checkShareCode(shareCode, request)
//     {
//         // check header first
//         if (request->header('share_code')) {
//             return checkCode(shareCode, request->header('share_code'));
//         }
//         // check query
//         if (request->input('share_code')) {
//             return checkCode(shareCode, request->input('share_code'));
//         }
//         return false;
//     }

//     func checkCode(valid, compare)
//     {
//         if (valid === compare) {
//             return true;
//         }
//         return false;
//     }

//     func deleteAccess(slug)
//     {
//         file = Access::where('slug', slug)->first();
//         if (deleteFile(file->file_id)) {
//             return file->delete();
//         };

//     }

//     func deleteFile(id)
//     {
//         file = File::findOrFail(id);
//         if (unlink(file->full_path)) {
//             return file->delete();
//         }
//     }

//     func reArrayFiles(&file_post)
//     {

//         file_ary = array();
//         file_count = count(file_post['name']);
//         file_keys = array_keys(file_post);

//         for (i = 0; i < file_count; i++) {
//             foreach (file_keys as key) {
//                 file_ary[i][key] = file_post[key][i];
//             }
//         }
//         return file_ary;
//     }

//     func generateSlug()
//     {
//         return generateRandomString(20);
//     }

//     func generateCode()
//     {
//         return generateRandomString(10);
//     }

//     func generateFilePath(mimeType, organisation, type, owner, date)
//     {
//         // check organisation
//         path = '/' . str_replace(' ', '-', organisation);
//         // move file
//         if (type) {
//             path = path . '/' . str_replace(' ', '-', type);
//             if (owner) {
//                 path = path . '/' . str_replace(' ', '-', owner);
//             }
//         } else {
//             path = path . '/' . mimeType;
//         }

//         if (date) {
//             path = path . '/' . date("Y/M/D", date);
//         } else {
//             path = path . '/' . date("Y") . '/' . date("M") . '/' . date("D");
//         }

//         if (env('MEDIA_OS') === 'windows') {
//             path = str_replace('/', '\\', path);
//         }

//         return env('MEDIA_PATH') . path;
//     }

//     func generateRandomString(length)
//     {
//         characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
//         charactersLength = strlen(characters);
//         randomString = '';
//         for (i = 0; i < length; i++) {
//             randomString .= characters[rand(0, charactersLength - 1)];
//         }
//         return randomString;
//     }

//     func fileBoolean(bool)
//     {
//         if(bool === 1 || bool === true || bool === '1' || bool === 'true' || bool === 0 || bool === false || bool === '0' || bool === 'false') {
//             return dbBool(bool);
//         }
//         return true;
//     }

//     func dbBool(bool)
//     {
//         if(env('DB_CONNECTION','mysql')) {
//             return convertBoolean(bool);
//         }
//         return checkBoolean(bool);
//     }

// }
// func createFile(filePath string) (*os.File, error) {
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return file, nil
// }

// func hashFileBlake2(filePath string) (string, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()

// 	hasher, err := blake2b.New256(nil) // You can choose different sizes like New512
// 	if err != nil {
// 		return "", err
// 	}

// 	if _, err := io.Copy(hasher, file); err != nil {
// 		return "", err
// 	}

// 	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
// }
