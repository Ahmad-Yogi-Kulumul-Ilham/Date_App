import 'package:geolocator/geolocator.dart';

final class Coordinates {
  const Coordinates(this.latitude, this.longitude);

  final double latitude;
  final double longitude;
}

enum LocationFailure {
  serviceDisabled,
  permissionDenied,
  permissionDeniedForever,
}

final class LocationException implements Exception {
  const LocationException(this.reason);

  final LocationFailure reason;

  String get message => switch (reason) {
    LocationFailure.serviceDisabled =>
      'Layanan lokasi di HP kamu sedang mati. Nyalakan dulu ya.',
    LocationFailure.permissionDenied =>
      'Izinkan akses lokasi supaya kami bisa menampilkan orang di sekitarmu.',
    LocationFailure.permissionDeniedForever =>
      'Akses lokasi ditolak permanen. Aktifkan lewat pengaturan aplikasi.',
  };
}

abstract interface class LocationService {
  /// Returns the current position, asking for permission if needed.
  /// Throws [LocationException] when the position is unavailable.
  Future<Coordinates> currentPosition();

  Future<void> openSettings();
}

final class GeolocatorLocationService implements LocationService {
  @override
  Future<Coordinates> currentPosition() async {
    if (!await Geolocator.isLocationServiceEnabled()) {
      throw const LocationException(LocationFailure.serviceDisabled);
    }

    var permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
    }
    switch (permission) {
      case LocationPermission.denied:
      case LocationPermission.unableToDetermine:
        throw const LocationException(LocationFailure.permissionDenied);
      case LocationPermission.deniedForever:
        throw const LocationException(LocationFailure.permissionDeniedForever);
      case LocationPermission.whileInUse:
      case LocationPermission.always:
        break;
    }

    // City-level accuracy is enough for "people nearby" and saves battery.
    final position = await Geolocator.getCurrentPosition(
      locationSettings: const LocationSettings(
        accuracy: LocationAccuracy.medium,
      ),
    );
    return Coordinates(position.latitude, position.longitude);
  }

  @override
  Future<void> openSettings() => Geolocator.openAppSettings();
}
