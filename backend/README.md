# Leo — backend API (Goravel)

Lihat [README di root](../README.md) untuk gambaran proyek dan daftar endpoint.

## Menjalankan

```bash
# dari root repo: nyalakan MySQL dan Redis
docker compose up -d

cd backend
cp .env.example .env
go run . artisan key:generate
go run . artisan jwt:secret
go run . artisan migrate
go run . artisan db:seed --seeder=DemoSeeder   # opsional: data demo di Bandung
go run .                                        # http://localhost:3000
```

## Struktur

```
app/
  models/          entitas + aturan kecil domain (umur, pasangan match)
  services/        logika bisnis; ports.go berisi interface yang dibutuhkan service
  repositories/    implementasi ports dengan ORM Goravel + MySQL
  http/
    controllers/   tipis: validasi -> panggil service -> kirim respons
    requests/      aturan validasi request
    resources/     bentuk JSON respons (model tidak pernah diserialisasi langsung)
    responses/     amplop {"data": ...} / {"error": {...}} dan pemetaan error
    middleware/    autentikasi JWT dan rate limiter
  apperr/          kode error domain yang stabil untuk klien
  geo/             bounding box dan pembulatan jarak
routes/
  api.go           daftar endpoint /api/v1
  dependencies.go  satu-satunya tempat merangkai repository -> service -> controller
database/
  migrations/      skema tabel
  seeders/         DemoSeeder
```

## Test

```bash
go test ./app/...
```

Service dites dengan repository palsu di memori (`app/services/fakes_test.go`), jadi test tidak butuh database.
