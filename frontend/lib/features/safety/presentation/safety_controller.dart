import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../matches/data/api_matches_repository.dart';
import '../../matches/presentation/matches_provider.dart';
import '../data/api_safety_repository.dart';
import '../domain/report_reason.dart';

final safetyControllerProvider =
    NotifierProvider.autoDispose<SafetyController, bool>(SafetyController.new);

/// Ending contact with someone: unmatching, blocking and reporting. The state
/// is whether an action is in flight. Each action can remove a match, so the
/// match list is refreshed after a success.
class SafetyController extends Notifier<bool> {
  @override
  bool build() => false;

  Future<void> unmatch(int matchId) =>
      _busy(() => ref.read(matchesRepositoryProvider).unmatch(matchId));

  Future<void> block(int userId) =>
      _busy(() => ref.read(safetyRepositoryProvider).block(userId));

  Future<void> report(int userId, ReportReason reason, {String details = ''}) =>
      _busy(
        () => ref
            .read(safetyRepositoryProvider)
            .report(userId, reason, details: details),
      );

  Future<void> _busy(Future<void> Function() action) async {
    state = true;
    try {
      await action();
      if (ref.mounted) ref.invalidate(matchesProvider);
    } finally {
      if (ref.mounted) state = false;
    }
  }
}
