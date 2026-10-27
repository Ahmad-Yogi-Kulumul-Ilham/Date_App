import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../profile/domain/user_profile.dart';
import '../data/api_auth_repository.dart';
import '../domain/auth_repository.dart';

/// The signed-in user, or null when signed out. The router redirects on changes.
///
/// Action methods throw on failure (for example wrong credentials) and leave the
/// state untouched, so forms can show the error without leaving the page.
final sessionControllerProvider =
    AsyncNotifierProvider<SessionController, UserProfile?>(
      SessionController.new,
    );

class SessionController extends AsyncNotifier<UserProfile?> {
  AuthRepository get _repository => ref.read(authRepositoryProvider);

  @override
  Future<UserProfile?> build() {
    final repository = ref.watch(authRepositoryProvider);
    final subscription = repository.sessionExpired.listen(
      (_) => state = const AsyncData(null),
    );
    ref.onDispose(subscription.cancel);
    return repository.restoreSession();
  }

  Future<void> login({required String email, required String password}) async {
    final user = await _repository.login(email: email, password: password);
    state = AsyncData(user);
  }

  Future<void> register(RegisterData data) async {
    final user = await _repository.register(data);
    state = AsyncData(user);
  }

  Future<void> logout() async {
    await _repository.logout();
    state = const AsyncData(null);
  }

  /// Replaces the cached profile after it was changed elsewhere (edit, photos).
  void updateProfile(UserProfile profile) => state = AsyncData(profile);
}
