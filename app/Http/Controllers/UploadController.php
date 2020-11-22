<?php

namespace App\Http\Controllers;

use App\Models\Access;
use App\Models\File;
use App\Helpers\Helper;
use Illuminate\Http\Request;

class UploadController extends Controller {
    // display uploader
    public function index() {
        return view('uploader');
    }

    public function store(Request $request) {
        $organisation = $request->organisation;
        $owner = $request->owner;
        $type = $request->type; 
        $public = $request->public;
        $date = $request->date; 
        if($request->file('media')) {
            
           switch($request->multi) {
                case true:
                return Helper::multiFileUpload($request->file('media'), $organisation, $owner, $type, $public, $date);

                default:
                return Helper::upload($request->file('media'), $organisation, $owner, $type, $public, $date);
            } 
        }
        return response()->json(['Error'=> 'No file/s found.'], 503);
        
    }

    public function view(Request $request, $slug) {
        $access = Access::where('slug', $slug)->first();
        if($access) {
            if(Helper::checkAccess($access, $request)) {
                $file = File::find($access->file_id);
                if (file_exists($file->full_path)) {
                    header("Content-Type: ".$file->type);
                    header("Content-Length:".$file->size);
                    header('Content-Disposition: inline; filename="' . $file->name . '"'); 
                    header('Content-Transfer-Encoding: binary'); 
                    header('Accept-Ranges: bytes'); 
                    $fp = fopen($file->full_path, 'rb');
                    return fpassthru($fp);      
                } else {
                    return response()->json(['Error'=> 'File not found.'], 503);
                } 
            }
            return response()->json(['error'=> 'Invalid Access Rights'],401);
        }
        return response()->json(['error' => 'Nothing found!'],404);
    }

    public function show(Request $request, $slug) {
        $access = Access::where('slug', $slug)->first();
        if($access) {
            if(Helper::checkAccessCode($access->access_code, $request)) {
                $file = File::find($access->file_id);
                if (file_exists($file->full_path)) {
                    return [
                        'access' => $access,
                        'file' => $file,
                    ];;      
                } else {
                    return response()->json(['Error'=> 'File not found.'], 503);
                } 
            }
            return response()->json(['error'=> 'Invalid Access Rights'],401);
        }
        return response()->json(['error' => 'Nothing found!'],404);
    }

    public function update(Request $request, $slug)
    {
        $access = Access::where('slug', $slug)->first();
        if($access) {
            if(Helper::checkAccessCode($access->access_code, $request)) {
                if(Helper::updateAccess($slug, $request)) {
                    return response()->json(['message' => 'File updated!'], 204);
                }
                return response()->json(['error' => 'An error Occured!'], 503);
            }
            return response()->json(['error'=> 'Invalid Access Rights'],401);
        }
        return response()->json(['error' => 'Nothing found!'],404);
    }

    public function destroy(Request $request, $slug)
    {
        $access = Access::where('slug', $slug)->first();
        if($access) {
            if(Helper::checkAccessCode($access->access_code, $request)) {
                if(Helper::deleteAccess($slug)) {
                    return response()->json(['message' => 'File deleted!'], 204);
                }
                return response()->json(['error' => 'An error Occured!'], 503);
            }
            return response()->json(['error'=> 'Invalid Access Rights'],401);
        }
        return response()->json(['error' => 'Nothing found!'],404);
        
    }
}