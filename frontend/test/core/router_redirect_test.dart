import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/router/app_router.dart';
import 'package:leo_chat/core/router/routes.dart';
import 'package:leo_chat/features/profile/domain/user_profile.dart';

import '../helpers/fakes.dart';

void main() {
  const loading = AsyncLoading<UserProfile?>();
  const signedOut = AsyncData<UserProfile?>(null);
  final signedIn = AsyncData<UserProfile?>(testProfile());

  test('waits on the splash while the session is restored', () {
    expect(redirectFor(loading, Routes.discover), Routes.splash);
    expect(redirectFor(loading, Routes.splash), isNull);
  });

  test(
    'stays on the splash when restoring failed, so it can offer a retry',
    () {
      final failed = AsyncError<UserProfile?>(
        Exception('offline'),
        StackTrace.empty,
      );
      expect(redirectFor(failed, Routes.discover), Routes.splash);
    },
  );

  test('sends signed-out users to login but lets them reach register', () {
    expect(redirectFor(signedOut, Routes.discover), Routes.login);
    expect(redirectFor(signedOut, Routes.chat(3)), Routes.login);
    expect(redirectFor(signedOut, Routes.register), isNull);
  });

  test('moves signed-in users off auth screens', () {
    expect(redirectFor(signedIn, Routes.login), Routes.discover);
    expect(redirectFor(signedIn, Routes.splash), Routes.discover);
    expect(redirectFor(signedIn, Routes.matches), isNull);
  });
}
