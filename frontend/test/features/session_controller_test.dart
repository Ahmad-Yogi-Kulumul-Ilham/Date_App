import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/network/api_exception.dart';
import 'package:leo_chat/features/auth/data/api_auth_repository.dart';
import 'package:leo_chat/features/auth/presentation/session_controller.dart';

import '../helpers/fakes.dart';

void main() {
  late FakeAuthRepository auth;
  late ProviderContainer container;

  setUp(() {
    auth = FakeAuthRepository();
    container = ProviderContainer(
      overrides: [authRepositoryProvider.overrideWithValue(auth)],
      retry: (_, _) => null,
    );
    addTearDown(container.dispose);
  });

  test('starts signed out when there is no stored session', () async {
    expect(await container.read(sessionControllerProvider.future), isNull);
  });

  test('restores a stored session', () async {
    auth.restored = testProfile(name: 'Ayu');

    final user = await container.read(sessionControllerProvider.future);

    expect(user?.name, 'Ayu');
  });

  test('login stores the user', () async {
    await container.read(sessionControllerProvider.future);

    await container
        .read(sessionControllerProvider.notifier)
        .login(email: 'rina@example.com', password: 'rahasia123');

    expect(container.read(sessionControllerProvider).value?.name, 'Rina');
  });

  test('failed login throws and keeps the user signed out', () async {
    await container.read(sessionControllerProvider.future);
    auth.loginError = const ApiException(
      code: 'invalid_credentials',
      message: 'Email atau password salah.',
      statusCode: 401,
    );

    await expectLater(
      container
          .read(sessionControllerProvider.notifier)
          .login(email: 'rina@example.com', password: 'salah'),
      throwsA(isA<ApiException>()),
    );
    expect(container.read(sessionControllerProvider).value, isNull);
  });

  test('an expired session signs the user out', () async {
    auth.restored = testProfile();
    await container.read(sessionControllerProvider.future);

    auth.expireSession();
    await Future<void>.delayed(Duration.zero);

    expect(container.read(sessionControllerProvider).value, isNull);
  });

  test('logout clears the user', () async {
    auth.restored = testProfile();
    await container.read(sessionControllerProvider.future);

    await container.read(sessionControllerProvider.notifier).logout();

    expect(auth.logoutCalls, 1);
    expect(container.read(sessionControllerProvider).value, isNull);
  });
}
