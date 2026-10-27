import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/network/api_exception.dart';
import '../../../core/providers.dart';
import '../../../core/storage/token_storage.dart';
import '../../profile/domain/user_profile.dart';
import '../domain/auth_repository.dart';

final authRepositoryProvider = Provider<AuthRepository>(
  (ref) => ApiAuthRepository(
    ref.watch(apiClientProvider),
    ref.watch(tokenStorageProvider),
  ),
);

final class ApiAuthRepository implements AuthRepository {
  const ApiAuthRepository(this._api, this._tokens);

  final ApiClient _api;
  final TokenStorage _tokens;

  @override
  Stream<void> get sessionExpired => _api.sessionExpired;

  @override
  Future<UserProfile?> restoreSession() async {
    if (await _tokens.read() == null) return null;
    try {
      final json = await _api.get<Map<String, dynamic>>('/me');
      return UserProfile.fromJson(json);
    } on ApiException catch (e) {
      if (e.isUnauthorized) return null;
      rethrow;
    }
  }

  @override
  Future<UserProfile> login({
    required String email,
    required String password,
  }) => _authenticate('/auth/login', {
    'email': email.trim(),
    'password': password,
  });

  @override
  Future<UserProfile> register(RegisterData data) =>
      _authenticate('/auth/register', data.toJson());

  @override
  Future<void> logout() async {
    try {
      await _api.post<void>('/auth/logout');
    } on ApiException {
      // The local session is cleared below either way.
    }
    await _tokens.clear();
  }

  Future<UserProfile> _authenticate(
    String path,
    Map<String, dynamic> body,
  ) async {
    final json = await _api.post<Map<String, dynamic>>(path, body: body);
    await _tokens.write(json['access_token'] as String);
    return UserProfile.fromJson(json['user'] as Map<String, dynamic>);
  }
}
