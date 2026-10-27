import '../../../core/location/location_service.dart';
import 'user_profile.dart';

abstract interface class ProfileRepository {
  /// Partial update: null arguments are left unchanged.
  Future<UserProfile> updateProfile({
    String? name,
    String? bio,
    Preference? interestedIn,
  });

  Future<void> updateLocation(Coordinates coordinates);

  Future<Photo> uploadPhoto(String filePath);

  Future<void> deletePhoto(int photoId);
}
