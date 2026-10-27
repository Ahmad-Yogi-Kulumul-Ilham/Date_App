import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/location/location_service.dart';
import '../../../core/network/api_client.dart';
import '../../../core/providers.dart';
import '../domain/profile_repository.dart';
import '../domain/user_profile.dart';

final profileRepositoryProvider = Provider<ProfileRepository>(
  (ref) => ApiProfileRepository(ref.watch(apiClientProvider)),
);

final class ApiProfileRepository implements ProfileRepository {
  const ApiProfileRepository(this._api);

  final ApiClient _api;

  @override
  Future<UserProfile> updateProfile({
    String? name,
    String? bio,
    Preference? interestedIn,
  }) async {
    final json = await _api.patch<Map<String, dynamic>>(
      '/me',
      body: {'name': ?name, 'bio': ?bio, 'interested_in': ?interestedIn?.name},
    );
    return UserProfile.fromJson(json);
  }

  @override
  Future<void> updateLocation(Coordinates coordinates) => _api.put<void>(
    '/me/location',
    body: {
      'latitude': coordinates.latitude,
      'longitude': coordinates.longitude,
    },
  );

  @override
  Future<Photo> uploadPhoto(String filePath) async {
    final form = FormData.fromMap({
      'photo': await MultipartFile.fromFile(filePath),
    });
    final json = await _api.post<Map<String, dynamic>>(
      '/me/photos',
      body: form,
    );
    return Photo.fromJson(json);
  }

  @override
  Future<void> deletePhoto(int photoId) => _api.delete('/me/photos/$photoId');
}
