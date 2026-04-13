<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Matched extends Model
{
    use HasFactory;

    protected $table = 'matcheds'; // Ensure correct table name

    protected $fillable = [
        'user_id_1', 'user_id_2'
    ];

    public function user1()
    {
        return $this->belongsTo(User::class, 'user_id_1');
    }

    public function user2()
    {
        return $this->belongsTo(User::class, 'user_id_2');
    }

    public function messages()
    {
        return $this->hasMany(Message::class, 'matched_id');
    }
}
