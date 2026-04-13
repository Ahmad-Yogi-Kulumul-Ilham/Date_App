<?php
namespace Database\Seeders;

use Illuminate\Database\Seeder;
use App\Models\User;
use App\Models\Profile;
use Illuminate\Support\Facades\Hash;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        $users = [
            ['name' => 'Alice Spectra', 'email' => 'alice@example.com'],
            ['name' => 'Elena Rosa', 'email' => 'elena@example.com'],
            ['name' => 'Sophia V', 'email' => 'sophia@example.com'],
        ];

        foreach($users as $u) {
            $user = User::create([
                'name' => $u['name'],
                'email' => $u['email'],
                'password' => Hash::make('password')
            ]);
            Profile::create([
                'user_id' => $user->id,
                'name' => $user->name,
                'bio' => 'Looking for someone to go on adventures with!',
                'photo_url' => 'https://ui-avatars.com/api/?name='.urlencode($user->name).'&size=400&background=random'
            ]);
        }
    }
}
