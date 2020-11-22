<?php

namespace App\Http\Controllers;

use App\Models\Access;
use App\Models\File;
use App\Helpers\Helper;
use Illuminate\Http\Request;

class AdminController extends Controller {
    // show files with access
    public function index(Request $request) {
        return Access::all();
    }
    // show public files
    public function public(Request $request) {
        return Access::where('public', 1)->get();
    }



}