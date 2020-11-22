<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Recovery extends Model 
{

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'email', 
        'ip', 
        'domain',
        'code', 
        'access_id'
    ];


}
