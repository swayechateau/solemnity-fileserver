<?php
namespace App\Helpers;

use App\Models\Access;
use App\Models\File;

class Helper
{
    public static function upload($files, $organisation, $owner, $type, $public, $date)
    {
        if (is_array($files)) {
            return Helper::multiFileUpload($files, $organisation, $owner, $type, $public, $date);
        }
        return Helper::uploadFile($files, $organisation, $owner, $type, $public, $date);
    }

    public static function multiFileUpload($files, $organisation, $owner, $type, $public, $date)
    {
        $stack = [];
        foreach ($files as $file) {
            array_push($stack, Helper::uploadFile($file, $organisation, $owner, $type, $public, $date));
        }
        return $stack;
    }

    public static function uploadFile($file, $organisation, $owner, $type, $public, $date)
    {

        if (empty($organisation)) {
            $organisation = 'global';
        }

        $public = Helper::fileBoolean($public);


        $destinationPath = Helper::generateFilePath($file->getMimeType(), $organisation, $type, $owner, $date);
        // check if file name exists
        $name = $file->getClientOriginalName();
        $unique = false;
        $looped = 1;
        if (file_exists($destinationPath . '/' . $name)) {
            while ($unique === false) {
                if (file_exists($destinationPath . '/' . '(' . $looped . ')_' . $name)) {
                    $looped++;
                } else {
                    $unique = true;
                }
            }
            $name = '(' . $looped . ')_' . $name;
        }

        // upload file
        $uploadedFile = File::create([
            "name" => $name,
            "extension" => $file->getClientOriginalExtension(),
            "path" => $destinationPath,
            "full_path" => $destinationPath . '/' . $name,
            "size" => $file->getSize(),
            "type" => $file->getMimeType(),
        ]);
        $file->move($destinationPath, $name);

        // create file access
        $slug = str_replace(' ', '-', $organisation) . Helper::generateSlug() . $uploadedFile->id;
        $shareCode = 's' . $uploadedFile->id . Helper::generateCode();
        $accessCode = 'a' . $uploadedFile->id . Helper::generateCode();
        $access = Access::create([
            'organisation' => $organisation,
            'owner' => $owner,
            'type' => $type,
            'public' => $public,
            'slug' => $slug,
            'share_code' => $shareCode,
            'access_code' => $accessCode,
            'file_id' => $uploadedFile->id,
        ]);
        // return file with access
        if ($access && $uploadedFile) {
            return [
                'access' => $access,
                'file' => $uploadedFile,
            ];
        }
        return false;
    }

    public static function updateAccess($slug, $request)
    {
        $path = $request->path;
        $name = $request->name;
        $access = Access::where('slug', $slug);
        if ($request->organisation) {
            $access->organisation = $request->organisation;
        }
        if ($request->owner) {
            $access->owner = $request->owner;
        }
        if ($request->type) {
            $access->type = $request->type;
        }
        if ($request->public) {
            $access->public = Helper::convertBoolean($request->public);
        }

        if ($path || $name) {
            $file = File::find($access->file_id);
            $updatedFile = Helper::updateFile($file, $path, $name);
            if ($updatedFile) {return true;}
        }

        return false;
    }

    public static function updateFile($file, $cpath, $cname)
    {
        if ($cname || $cpath) {
            $ext = $file->extension;
            $path = $file->path;
            $name = $file->name;

            // change path
            if ($cpath) {
                $path = $file->path;
                if (is_writable($path)) {
                    $name = Helper::checkFileName($path, $name);
                    $file->path = $path;
                }
                return false;
            }

            // check if name has .extention
            // add .extention to name end if missing
            if ($cname) {
                $name = $cname;
                if (!strpos($name, $ext)) {
                    $name = $name . $ext;
                }
                $name = Helper::checkFileName($path, $name);
                $file->name = $name;
            }
            // update file
            rename($file->full_path, $path . '/' . $name);
            // update full_path
            $file->full_path = $path . '/' . $name;
            //update file record
            return $file->save();
        }

        return false;

    }

    public static function checkFileName($path, $name)
    {
        $unique = false;
        $looped = 1;
        while ($unique === false) {
            if (file_exists($path . '/' . $name)) {
                $name = '(' . $looped . ')_' . $name;
                $looped++;
            } else {
                $unique = true;
            }
        }
        return $name;
    }

    public static function checkAccess($access, $request)
    {

        if (Helper::checkBoolean($access->public) === false) {
            if (Helper::checkAccessCode($access->access_code, $request)) {
                return true;
            }

            if (Helper::checkShareCode($access->share_code, $request)) {
                return true;
            }

            return false;
        }

        return true;
    }

    public static function checkAccessCode($accessCode, $request)
    {

        if ($request->header('access_code')) {
            return Helper::checkCode($accessCode, $request->header('access_code'));
        }

        if ($request->input('access_code')) {
            return Helper::checkCode($accessCode, $request->input('access_code'));
        }

        return false;

    }

    public static function checkShareCode($shareCode, $request)
    {
        // check header first
        if ($request->header('share_code')) {
            return Helper::checkCode($shareCode, $request->header('share_code'));
        }
        // check query
        if ($request->input('share_code')) {
            return Helper::checkCode($shareCode, $request->input('share_code'));
        }
        return false;
    }

    public static function checkCode($valid, $compare)
    {
        if ($valid === $compare) {
            return true;
        }
        return false;
    }

    public static function deleteAccess($slug)
    {
        $file = Access::where('slug', $slug)->first();
        if (Helper::deleteFile($file->file_id)) {
            return $file->delete();
        };

    }

    public static function deleteFile($id)
    {
        $file = File::findOrFail($id);
        if (unlink($file->full_path)) {
            return $file->delete();
        }
    }

    public static function reArrayFiles(&$file_post)
    {

        $file_ary = array();
        $file_count = count($file_post['name']);
        $file_keys = array_keys($file_post);

        for ($i = 0; $i < $file_count; $i++) {
            foreach ($file_keys as $key) {
                $file_ary[$i][$key] = $file_post[$key][$i];
            }
        }
        return $file_ary;
    }

    public static function generateSlug()
    {
        return Helper::generateRandomString(20);
    }

    public static function generateCode()
    {
        return Helper::generateRandomString(10);
    }

    public static function generateFilePath($mimeType, $organisation, $type, $owner, $date)
    {
        // check organisation
        $path = '/' . str_replace(' ', '-', $organisation);
        // move file
        if ($type) {
            $path = $path . '/' . str_replace(' ', '-', $type);
            if ($owner) {
                $path = $path . '/' . str_replace(' ', '-', $owner);
            }
        } else {
            $path = $path . '/' . $mimeType;
        }

        if ($date) {
            $path = $path . '/' . date("Y/M/D", $date);
        } else {
            $path = $path . '/' . date("Y") . '/' . date("M") . '/' . date("D");
        }

        if (env('MEDIA_OS') === 'windows') {
            $path = str_replace('/', '\\', $path);
        }

        return env('MEDIA_PATH') . $path;
    }

    public static function generateRandomString($length)
    {
        $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        $charactersLength = strlen($characters);
        $randomString = '';
        for ($i = 0; $i < $length; $i++) {
            $randomString .= $characters[rand(0, $charactersLength - 1)];
        }
        return $randomString;
    }

    public static function fileBoolean($bool)
    {
        if($bool === 1 || $bool === true || $bool === '1' || $bool === 'true' || $bool === 0 || $bool === false || $bool === '0' || $bool === 'false') {
            return Helper::dbBool($bool);
        }
        return true;
    }

    public static function dbBool($bool)
    {
        if(env('DB_CONNECTION','mysql')) {
            return Helper::convertBoolean($bool);
        }
        return Helper::checkBoolean($bool);
    }
    public static function checkBoolean($bool)
    {
        if ($bool === 1 || $bool === true || $bool === '1' || $bool === 'true') {
            return true;
        }
        return false;
    }

    public static function convertBoolean($bool)
    {
        if ($bool === 1 || $bool === true || $bool === '1' || $bool === 'true') {
            return 1;
        }
        return 0;
    }
}
