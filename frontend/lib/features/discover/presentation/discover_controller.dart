import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_exception.dart';
import '../../../core/providers.dart';
import '../../auth/presentation/current_user.dart';
import '../../matches/presentation/matches_provider.dart';
import '../../profile/data/api_profile_repository.dart';
import '../data/api_discover_repository.dart';
import '../domain/candidate.dart';

final discoverControllerProvider =
    AsyncNotifierProvider.autoDispose<DiscoverController, List<Candidate>>(
      DiscoverController.new,
    );

/// The stack of people to swipe on. Loading it first refreshes the user's
/// location so results are always relative to where they are now.
class DiscoverController extends AsyncNotifier<List<Candidate>> {
  @override
  Future<List<Candidate>> build() async {
    ref.watch(currentUserIdProvider);

    final position = await ref.read(locationServiceProvider).currentPosition();
    await ref.read(profileRepositoryProvider).updateLocation(position);
    return ref.read(discoverRepositoryProvider).fetchCandidates();
  }

  /// Removes [candidate] from the stack immediately and records the swipe.
  /// If the request fails the card comes back and the error is rethrown.
  Future<SwipeOutcome> swipe(
    Candidate candidate,
    SwipeDirection direction,
  ) async {
    state = AsyncData([
      for (final c in state.value ?? const <Candidate>[])
        if (c.id != candidate.id) c,
    ]);

    try {
      final outcome = await ref
          .read(discoverRepositoryProvider)
          .swipe(candidate.id, direction);
      if (outcome.matched) ref.invalidate(matchesProvider);
      return outcome;
    } on ApiException catch (e) {
      if (e.code == 'already_swiped') return const SwipeOutcome(matched: false);
      if (ref.mounted) state = AsyncData([candidate, ...?state.value]);
      rethrow;
    }
  }
}
