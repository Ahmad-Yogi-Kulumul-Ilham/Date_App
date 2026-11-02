import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/network/api_exception.dart';
import 'package:leo_chat/features/auth/data/api_auth_repository.dart';
import 'package:leo_chat/features/auth/presentation/session_controller.dart';
import 'package:leo_chat/features/matches/data/api_matches_repository.dart';
import 'package:leo_chat/features/matches/presentation/matches_provider.dart';
import 'package:leo_chat/features/safety/data/api_safety_repository.dart';
import 'package:leo_chat/features/safety/domain/report_reason.dart';
import 'package:leo_chat/features/safety/presentation/safety_controller.dart';

import '../helpers/fakes.dart';

void main() {
  late FakeMatchesRepository matches;
  late FakeSafetyRepository safety;
  late ProviderContainer container;

  setUp(() {
    matches = FakeMatchesRepository();
    safety = FakeSafetyRepository();
    container = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(
          FakeAuthRepository(restored: testProfile()),
        ),
        matchesRepositoryProvider.overrideWithValue(matches),
        safetyRepositoryProvider.overrideWithValue(safety),
      ],
      retry: (_, _) => null,
    );
    addTearDown(container.dispose);
    container.listen(safetyControllerProvider, (_, _) {});
  });

  SafetyController controller() =>
      container.read(safetyControllerProvider.notifier);

  Future<void> loadMatches() async {
    await container.read(sessionControllerProvider.future);
    container.listen(matchesProvider, (_, _) {});
    await container.read(matchesProvider.future);
  }

  test('unmatching refreshes the match list', () async {
    await loadMatches();
    final fetchesBefore = matches.fetchCalls;

    await controller().unmatch(5);
    await container.read(matchesProvider.future);

    expect(matches.unmatched, [5]);
    expect(matches.fetchCalls, fetchesBefore + 1);
  });

  test('reporting sends the reason and details', () async {
    await controller().report(2, ReportReason.spam, details: 'jualan pulsa');

    expect(safety.reports, [(2, ReportReason.spam, 'jualan pulsa')]);
    expect(container.read(safetyControllerProvider), isFalse);
  });

  test('a failed block rethrows and clears the busy flag', () async {
    await loadMatches();
    final fetchesBefore = matches.fetchCalls;
    safety.error = const ApiException(code: 'network', message: 'offline');

    await expectLater(controller().block(2), throwsA(isA<ApiException>()));

    expect(container.read(safetyControllerProvider), isFalse);
    expect(safety.blocked, isEmpty);
    expect(matches.fetchCalls, fetchesBefore, reason: 'nothing changed');
  });
}
