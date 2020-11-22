<?php

namespace App\Http\Controllers;

use App\Models\Access;
use App\Models\File;
use App\Helpers\Helper;
use Illuminate\Support\Str;
use App\Http\Requests;
use Illuminate\Http\Request;

class DocsController extends Controller {

    public function index() {
        // display uploader
        return view('docs');
    }

    public function test(Request $request) {
        return $uuid = Str::uuid()->toString();
    }

}