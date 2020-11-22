<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Access extends Model 
{

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'organisation','owner','type',
        'public','slug','share_code',
        'access_code','file_id'
    ];


}
