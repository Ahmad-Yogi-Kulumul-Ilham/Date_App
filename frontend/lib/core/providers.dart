import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'location/location_service.dart';
import 'network/api_client.dart';
import 'storage/token_storage.dart';

/// App-wide infrastructure. Tests override these with fakes.

final tokenStorageProvider = Provider<TokenStorage>(
  (ref) => SecureTokenStorage(),
);

final apiClientProvider = Provider<ApiClient>((ref) {
  final client = ApiClient(tokens: ref.watch(tokenStorageProvider));
  ref.onDispose(client.dispose);
  return client;
});

final locationServiceProvider = Provider<LocationService>(
  (ref) => GeolocatorLocationService(),
);
