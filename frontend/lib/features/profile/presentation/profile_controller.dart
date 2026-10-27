import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../auth/presentation/session_controller.dart';
import '../data/api_profile_repository.dart';
import '../domain/profile_repository.dart';
import '../domain/user_profile.dart';

final profileControllerProvider =
    NotifierProvider.autoDispose<ProfileController, bool>(
      ProfileController.new,
    );

/// Profile edits. The state is whether a change is in flight; the profile
/// itself lives in [sessionControllerProvider] so every screen stays in sync.
class ProfileController extends Notifier<bool> {
  @override
  bool build() => false;

  ProfileRepository get _repository => ref.read(profileRepositoryProvider);
  SessionController get _session =>
      ref.read(sessionControllerProvider.notifier);
  UserProfile? get _current => ref.read(sessionControllerProvider).value;

  Future<void> updateDetails({
    String? name,
    String? bio,
    Preference? interestedIn,
  }) => _busy(() async {
    final updated = await _repository.updateProfile(
      name: name,
      bio: bio,
      interestedIn: interestedIn,
    );
    _session.updateProfile(updated);
  });

  Future<void> addPhoto(String filePath) => _busy(() async {
    final photo = await _repository.uploadPhoto(filePath);
    final current = _current;
    if (current != null) {
      _session.updateProfile(
        current.copyWith(photos: [...current.photos, photo]),
      );
    }
  });

  Future<void> deletePhoto(Photo photo) => _busy(() async {
    await _repository.deletePhoto(photo.id);
    final current = _current;
    if (current != null) {
      _session.updateProfile(
        current.copyWith(
          photos: [
            for (final p in current.photos)
              if (p.id != photo.id) p,
          ],
        ),
      );
    }
  });

  Future<void> _busy(Future<void> Function() action) async {
    state = true;
    try {
      await action();
    } finally {
      if (ref.mounted) state = false;
    }
  }
}
