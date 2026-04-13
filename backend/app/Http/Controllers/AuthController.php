<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\User;
use App\Models\Profile;
use Illuminate\Support\Facades\Hash;

class AuthController extends Controller
{
    public function register(Request $request) {
        $request->validate([
            'name' => 'required|string',
            'email' => 'required|email|unique:users',
            'password' => 'required|string|min:6'
        ]);

        $user = User::create([
            'name' => $request->name,
            'email' => $request->email,
            'password' => Hash::make($request->password)
        ]);

        Profile::create(['user_id' => $user->id, 'name' => $user->name, 'photo_url' => 'https://ui-avatars.com/api/?name='.urlencode($user->name)]);

        return response()->json(['token' => $user->createToken('auth')->plainTextToken]);
    }

    public function login(Request $request) {
        $request->validate(['email' => 'required|email', 'password' => 'required']);
        $user = User::where('email', $request->email)->first();

        if (!$user || !Hash::check($request->password, $user->password)) {
            return response()->json(['message' => 'Invalid credentials'], 401);
        }

        return response()->json(['token' => $user->createToken('auth')->plainTextToken, 'user' => $user->load('profile')]);
    }

    public function logout(Request $request) {
        $request->user()->currentAccessToken()->delete();
        return response()->json(['message' => 'Logged out']);
    }

    public function getProfile(Request $request) {
        return response()->json($request->user()->load('profile'));
    }

    public function updateProfile(Request $request) {
        $profile = $request->user()->profile;
        $profile->update($request->only('bio', 'gender', 'interested_in', 'birthdate'));
        return response()->json($profile);
    }
}
