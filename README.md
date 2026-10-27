# Leo

Aplikasi dating: temukan orang di sekitar, swipe, match, lalu chat.

| Bagian | Teknologi |
|---|---|
| Aplikasi mobile | Flutter, Riverpod 3, go_router, Dio |
| Backend API | Go + [Goravel](https://www.goravel.dev) v1.18 |
| Database | MySQL 8.4 |
| Cache, rate limit, blacklist token | Redis 7 |

```
.
├── backend/             API Goravel (lihat backend/README.md)
├── frontend/            aplikasi Flutter (lihat frontend/README.md)
├── docker-compose.yml   MySQL + Redis untuk development
└── .github/workflows/   CI: format, lint, dan test untuk kedua sisi
```

## Menjalankan secara lokal

Prasyarat: Go 1.25+, Flutter 3.47+, Docker Desktop.

```bash
# 1. Database dan Redis
docker compose up -d

# 2. Backend
cd backend
cp .env.example .env
go run . artisan key:generate
go run . artisan jwt:secret
go run . artisan migrate
go run . artisan db:seed --seeder=DemoSeeder
go run .

# 3. Aplikasi (terminal lain, dengan Android emulator menyala)
cd frontend
flutter pub get
flutter run
```

Akun demo dari seeder: **demo@leochat.id** / **password123**. Akun ini berlokasi di Bandung dengan 10 profil di sekitarnya. Rina, Nadia, dan Bima sudah menyukai akun demo, jadi like ke mereka langsung menghasilkan match.

Discover memakai lokasi HP. Di emulator, atur lokasi ke sekitar Bandung (-6.90, 107.61) lewat *Extended controls → Location* supaya profil demo muncul.

Foto profil disajikan dari URL berdasarkan `APP_URL`. Untuk Android emulator, isi `APP_URL=http://10.0.2.2:3000` di `backend/.env` supaya gambar bisa dimuat.

## API

Base URL `/api/v1`. Semua respons sukses berbentuk `{"data": ...}`, dan error berbentuk:

```json
{"error": {"code": "validation_failed", "message": "...", "fields": {"email": "..."}}}
```

`code` bersifat stabil sehingga aplikasi bisa bercabang berdasarkan kode ini. `message` sudah dalam Bahasa Indonesia dan siap ditampilkan ke pengguna.

| Method | Path | Keterangan |
|---|---|---|
| POST | `/auth/register` | Daftar dan langsung mendapat token (umur minimal 18) |
| POST | `/auth/login` | Login |
| POST | `/auth/refresh` | Tukar token kedaluwarsa dengan token baru |
| POST | `/auth/logout` | Masukkan token ke blacklist |
| GET | `/me` | Profil sendiri |
| PATCH | `/me` | Ubah nama, bio, atau preferensi |
| PUT | `/me/location` | Perbarui lokasi |
| POST | `/me/photos` | Unggah foto (multipart `photo`, jpg/png/webp, maksimal 5 MB, maksimal 6 foto) |
| DELETE | `/me/photos/{id}` | Hapus foto |
| GET | `/discover` | Kandidat terdekat yang belum di-swipe. Query opsional: `radius_km`, `min_age`, `max_age`, `limit` |
| POST | `/swipes` | `{"user_id": 2, "direction": "like" \| "pass"}` mengembalikan `{"matched": true, "match_id": 5}` |
| GET | `/matches` | Daftar match beserta pesan terakhir |
| GET | `/matches/{id}/messages` | Riwayat pesan, terbaru dulu. Pakai `before_id` untuk halaman berikutnya |
| POST | `/matches/{id}/messages` | Kirim pesan |

Semua endpoint selain register, login, dan refresh butuh header `Authorization: Bearer <token>`. Login dan register dibatasi 10 permintaan per menit per IP, swipe 120 per menit, dan chat 60 per menit per pengguna.

## Keputusan desain

- **Privasi lokasi.** Koordinat pengguna tidak pernah dikirim ke pengguna lain. API hanya memberi jarak yang dibulatkan ke kilometer (minimal "1 km"), supaya posisi seseorang tidak bisa ditriangulasi.
- **Cocok dua arah.** Discover hanya menampilkan orang yang gendernya sesuai preferensimu *dan* yang preferensinya sesuai gendermu.
- **Pencarian terdekat di MySQL.** Kotak batas (bounding box) di kolom `latitude` dan `longitude` yang berindeks menyaring kandidat dulu, lalu `ST_Distance_Sphere` menghitung jarak sebenarnya.
- **Match tanpa duplikat.** Pasangan disimpan berurutan (`user_one_id < user_two_id`) dengan unique index. Like bersamaan dari dua sisi tetap menghasilkan satu match.
- **ID match tidak bisa ditebak-tebak.** Orang luar yang mengakses match orang lain mendapat 404, bukan 403.
- **Lapisan yang bisa dites.** Service hanya bergantung pada interface (`backend/app/services/ports.go`), dan controller Flutter hanya bergantung pada repository abstrak. Keduanya dites dengan implementasi palsu, tanpa database atau server.

## Langkah berikutnya

- Chat realtime dengan WebSocket, ditambah Redis Pub/Sub saat server lebih dari satu. Saat ini chat memakai polling tiap 4 detik, dan perubahannya cukup di `ChatController`.
- Push notification (FCM) untuk match dan pesan baru.
- Fitur block dan report, plus moderasi foto.
- Simpan foto di object storage (Cloudflare R2 atau S3) lewat driver filesystem Goravel.
- Test integrasi repository terhadap MySQL sungguhan, misalnya dengan `facades.Testing().Docker()`.
