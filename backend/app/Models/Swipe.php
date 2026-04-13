<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Swipe extends Model
{
    use HasFactory;

    protected $fillable = [
        'swiper_id', 'swipee_id', 'type'
    ];

    public function swiper()
    {
        return $this->belongsTo(User::class, 'swiper_id');
    }

    public function swipee()
    {
        return $this->belongsTo(User::class, 'swipee_id');
    }
}
