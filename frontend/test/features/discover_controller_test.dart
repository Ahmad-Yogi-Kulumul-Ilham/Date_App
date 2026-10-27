import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/location/location_service.dart';
import 'package:leo_chat/core/network/api_exception.dart';
import 'package:leo_chat/core/providers.dart';
import 'package:leo_chat/features/auth/data/api_auth_repository.dart';
import 'package:leo_chat/features/discover/data/api_discover_repository.dart';
import 'package:leo_chat/features/discover/domain/candidate.dart';
import 'package:leo_chat/features/discover/presentation/discover_controller.dart';
import 'package:leo_chat/features/profile/data/api_profile_repository.dart';

import '../helpers/fakes.dart';

void main() {
  late FakeDiscoverRepository discover;
  late FakeProfileRepository profile;
  late FakeLocationService location;
  late ProviderContainer container;

  setUp(() {
    discover = FakeDiscoverRepository([
      testCandidate(2),
      testCandidate(3, name: 'Raka'),
    ]);
    profile = FakeProfileRepository();
    location = FakeLocationService();
    container = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(
          FakeAuthRepository(restored: testProfile()),
        ),
        discoverRepositoryProvider.overrideWithValue(discover),
        profileRepositoryProvider.overrideWithValue(profile),
        locationServiceProvider.overrideWithValue(location),
      ],
      retry: (_, _) => null,
    );
    addTearDown(container.dispose);
  });

  Future<List<Candidate>> load() {
    container.listen(discoverControllerProvider, (_, _) {});
    return container.read(discoverControllerProvider.future);
  }

  test('sends the current location before loading candidates', () async {
    final candidates = await load();

    expect(profile.lastLocation?.latitude, -6.9);
    expect(candidates.map((c) => c.id), [2, 3]);
  });

  test('surfaces location problems as a LocationException', () async {
    location.failure = LocationFailure.permissionDenied;

    await expectLater(load(), throwsA(isA<LocationException>()));
  });

  test('swiping removes the card and reports a mutual like', () async {
    final candidates = await load();
    discover.likesBack.add(2);

    final outcome = await container
        .read(discoverControllerProvider.notifier)
        .swipe(candidates.first, SwipeDirection.like);

    expect(outcome.matched, isTrue);
    expect(outcome.matchId, 99);
    expect(container.read(discoverControllerProvider).value?.map((c) => c.id), [
      3,
    ]);
  });

  test('a failed swipe puts the card back on top', () async {
    final candidates = await load();
    discover.swipeError = const ApiException(
      code: 'network',
      message: 'offline',
    );

    await expectLater(
      container
          .read(discoverControllerProvider.notifier)
          .swipe(candidates.first, SwipeDirection.pass),
      throwsA(isA<ApiException>()),
    );
    expect(container.read(discoverControllerProvider).value?.map((c) => c.id), [
      2,
      3,
    ]);
  });
}
