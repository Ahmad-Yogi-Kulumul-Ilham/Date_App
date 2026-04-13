<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Models\User;
use App\Models\Swipe;
use App\Models\Matched;

class SwipeController extends Controller
{
    public function discover(Request $request) {
        $user = $request->user();
        $swipedIds = Swipe::where('swiper_id', $user->id)->pluck('swipee_id')->toArray();
        $swipedIds[] = $user->id; // exclude self

        // Ideally we would filter by gender/interest, but keep simple for now
        $profiles = User::whereNotIn('id', $swipedIds)->with('profile')->inRandomOrder()->limit(10)->get();

        return response()->json($profiles);
    }

    public function swipe(Request $request) {
        $request->validate([
            'swipee_id' => 'required|exists:users,id',
            'type' => 'required|in:like,pass'
        ]);

        $user = $request->user();
        
        $swipe = Swipe::updateOrCreate(
            ['swiper_id' => $user->id, 'swipee_id' => $request->swipee_id],
            ['type' => $request->type]
        );

        $isMatch = false;
        if ($request->type === 'like') {
            // Auto match for testing
            Matched::firstOrCreate([
                'user_id_1' => min($user->id, $request->swipee_id),
                'user_id_2' => max($user->id, $request->swipee_id)
            ]);
            $isMatch = true;
        }

        return response()->json(['success' => true, 'match' => $isMatch]);
    }
}
