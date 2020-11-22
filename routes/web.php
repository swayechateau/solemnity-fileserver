<?php

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It is a breeze. Simply tell Lumen the URIs it should respond to
| and give it the Closure to call when that URI is requested.
|
*/

// view API Doc
$router->get('/', 'DocsController@index');

// Upload File
$router->get('upload','UploadController@index');
$router->post('upload','UploadController@store');

// View File
$router->get('view/{slug}','UploadController@view');

// View File Details
$router->get('file/{slug}','UploadController@show');
// Modify File
$router->put('file/{slug}','UploadController@update'); // needs some work

// Remove File
$router->delete('file/{slug}','UploadController@destroy');

// admin
$router->get('admin/files','AdminController@index');

// public global files
$router->get('public','AdminController@public');
