<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateAccessesTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
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
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('accesses');
    }
}
