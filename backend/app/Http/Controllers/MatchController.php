<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\Matched;
use App\Models\Message;

class MatchController extends Controller
{
    public function index(Request $request) {
        $user = $request->user();
        $matches = Matched::where('user_id_1', $user->id)
            ->orWhere('user_id_2', $user->id)
            ->with(['user1.profile', 'user2.profile'])
            ->get();
            
        return response()->json($matches);
    }

    public function getMessages(Request $request, Matched $matched) {
        $user = $request->user();
        if ($matched->user_id_1 !== $user->id && $matched->user_id_2 !== $user->id) {
            return response()->json(['message' => 'Unauthorized'], 403);
        }

        $messages = $matched->messages()->orderBy('created_at', 'asc')->get();
        return response()->json($messages);
    }

    public function sendMessage(Request $request, Matched $matched) {
        $request->validate(['content' => 'required|string']);
        $user = $request->user();
        
        if ($matched->user_id_1 !== $user->id && $matched->user_id_2 !== $user->id) {
            return response()->json(['message' => 'Unauthorized'], 403);
        }

        $message = Message::create([
            'matched_id' => $matched->id,
            'sender_id' => $user->id,
            'content' => $request->content
        ]);

        return response()->json($message);
    }
}
