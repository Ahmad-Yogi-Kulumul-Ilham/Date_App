# Leo — aplikasi Flutter

Lihat [README di root](../README.md) untuk gambaran proyek dan cara menjalankan backend.

## Menjalankan

```bash
flutter pub get

# Android emulator (default API: http://10.0.2.2:3000/api/v1)
flutter run

# iOS simulator
flutter run --dart-define=API_BASE_URL=http://localhost:3000/api/v1

# HP fisik di jaringan yang sama (ganti IP dengan IP laptop)
flutter run --dart-define=API_BASE_URL=http://192.168.1.10:3000/api/v1
```

## Struktur

```
lib/
  core/                 infrastruktur bersama: config, tema, network, router, widget umum
  features/<fitur>/
    domain/             model + interface repository (tanpa Flutter/Dio)
    data/               implementasi repository memakai ApiClient
    presentation/       halaman, widget, dan controller Riverpod
```

Aturannya: `presentation` hanya bicara ke `domain`, dan implementasi `data` disuntikkan lewat provider Riverpod. Karena itu setiap controller bisa dites dengan repository palsu (lihat `test/helpers/fakes.dart`).

## Cek kualitas

```bash
dart format lib test
flutter analyze
flutter test
```
