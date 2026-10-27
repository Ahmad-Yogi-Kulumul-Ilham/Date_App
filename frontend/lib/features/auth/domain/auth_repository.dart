import '../../profile/domain/user_profile.dart';

final class RegisterData {
  const RegisterData({
    required this.name,
    required this.email,
    required this.password,
    required this.birthDate,
    required this.gender,
    required this.interestedIn,
  });

  final String name;
  final String email;
  final String password;
  final DateTime birthDate;
  final Gender gender;
  final Preference interestedIn;

  Map<String, dynamic> toJson() => {
    'name': name.trim(),
    'email': email.trim(),
    'password': password,
    'birth_date': _formatDate(birthDate),
    'gender': gender.name,
    'interested_in': interestedIn.name,
  };

  static String _formatDate(DateTime d) =>
      '${d.year.toString().padLeft(4, '0')}-'
      '${d.month.toString().padLeft(2, '0')}-'
      '${d.day.toString().padLeft(2, '0')}';
}

abstract interface class AuthRepository {
  /// Returns the stored session's user, or null when nobody is signed in or
  /// the stored token is no longer valid.
  Future<UserProfile?> restoreSession();

  Future<UserProfile> login({required String email, required String password});

  Future<UserProfile> register(RegisterData data);

  Future<void> logout();

  /// Emits when the server rejects the session and it cannot be refreshed.
  Stream<void> get sessionExpired;
}
