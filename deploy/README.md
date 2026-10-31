# Deploy ke VPS

Target: `l.duzi.live` di VPS Ubuntu yang sudah menjalankan Nginx, Docker, dan MySQL. Leo hanya menambah dua container (backend dan Redis) dan memakai MySQL yang sudah ada di host, jadi tambahan RAM-nya sekitar 100 MB.

```
Internet → Nginx :443 (l.duzi.live) → backend 127.0.0.1:3003 → MySQL host 127.0.0.1:3306
                                                             → Redis      127.0.0.1:6380
```

Tidak ada port Leo yang terbuka ke publik. Semua lalu lintas masuk lewat Nginx.

## Pemasangan pertama

Semua perintah dijalankan sebagai root di VPS.

### 1. Ambil kode

```bash
git clone https://github.com/Ahmad-Yogi-Kulumul-Ilham/Date_App.git /opt/leochat
cd /opt/leochat
```

### 2. Database dan `.env`

```bash
DB_PASS=$(openssl rand -hex 20)

mysql -e "CREATE DATABASE IF NOT EXISTS leochat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS 'leochat'@'127.0.0.1' IDENTIFIED BY '$DB_PASS';
GRANT ALL PRIVILEGES ON leochat.* TO 'leochat'@'127.0.0.1';"

cp deploy/production.env.example backend/.env
sed -i "s/^APP_KEY=.*/APP_KEY=$(openssl rand -hex 16)/; s/^JWT_SECRET=.*/JWT_SECRET=$(openssl rand -hex 32)/; s/^DB_PASSWORD=.*/DB_PASSWORD=$DB_PASS/" backend/.env
chmod 600 backend/.env
```

Kalau `mysql` meminta password root, ganti dengan `mysql -u root -p -e "..."`.

### 3. Jalankan backend

```bash
docker compose -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.prod.yml exec backend ./main artisan migrate
curl http://127.0.0.1:3003/health
```

Akun demo bersifat opsional. Password-nya tercantum di README publik, jadi siapa pun bisa login dengan akun ini:

```bash
docker compose -f docker-compose.prod.yml exec backend ./main artisan db:seed --seeder=DemoSeeder --force
```

### 4. Nginx dan HTTPS

```bash
cp deploy/nginx/l.duzi.live.conf /etc/nginx/sites-available/l.duzi.live
ln -s /etc/nginx/sites-available/l.duzi.live /etc/nginx/sites-enabled/l.duzi.live
nginx -t && systemctl reload nginx
certbot --nginx -d l.duzi.live
curl https://l.duzi.live/health
```

## Update

```bash
cd /opt/leochat
git pull
docker compose -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.prod.yml exec backend ./main artisan migrate
```

## Aplikasi Android

Di laptop:

```bash
cd frontend
flutter build apk --dart-define=API_BASE_URL=https://l.duzi.live/api/v1
```

Hasilnya ada di `frontend/build/app/outputs/flutter-apk/app-release.apk`.

## Perawatan

```bash
docker compose -f docker-compose.prod.yml ps                # status
docker compose -f docker-compose.prod.yml logs -f backend   # log API
docker compose -f docker-compose.prod.yml restart backend
```

Foto profil dan log tersimpan di volume `leochat_storage`, data Redis di `leochat_redis-data`. Data utama ada di database `leochat` pada MySQL host, jadi pastikan database ini ikut dalam backup server.
