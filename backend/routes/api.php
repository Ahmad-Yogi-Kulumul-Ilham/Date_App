<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\AuthController;
use App\Http\Controllers\SwipeController;
use App\Http\Controllers\MatchController;

Route::post('/register', [AuthController::class, 'register']);
Route::post('/login', [AuthController::class, 'login']);

Route::middleware('auth:sanctum')->group(function () {
    Route::post('/logout', [AuthController::class, 'logout']);
    Route::get('/profile', [AuthController::class, 'getProfile']);
    Route::put('/profile', [AuthController::class, 'updateProfile']);

    Route::get('/discover', [SwipeController::class, 'discover']);
    Route::post('/swipe', [SwipeController::class, 'swipe']);

    Route::get('/matches', [MatchController::class, 'index']);
    Route::get('/matches/{matched}/messages', [MatchController::class, 'getMessages']);
    Route::post('/matches/{matched}/messages', [MatchController::class, 'sendMessage']);
});
