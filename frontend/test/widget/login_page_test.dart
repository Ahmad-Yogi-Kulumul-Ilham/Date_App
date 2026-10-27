import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/network/api_exception.dart';
import 'package:leo_chat/core/theme/app_theme.dart';
import 'package:leo_chat/features/auth/data/api_auth_repository.dart';
import 'package:leo_chat/features/auth/presentation/login_page.dart';
import 'package:leo_chat/features/auth/presentation/session_controller.dart';

import '../helpers/fakes.dart';

void main() {
  Future<ProviderContainer> pumpLogin(
    WidgetTester tester,
    FakeAuthRepository auth,
  ) async {
    final container = ProviderContainer(
      overrides: [authRepositoryProvider.overrideWithValue(auth)],
      retry: (_, _) => null,
    );
    addTearDown(container.dispose);

    await tester.pumpWidget(
      UncontrolledProviderScope(
        container: container,
        child: MaterialApp(
          theme: AppTheme.light(useGoogleFonts: false),
          home: const LoginPage(),
        ),
      ),
    );
    await container.read(sessionControllerProvider.future);
    return container;
  }

  testWidgets('validates the form before calling the API', (tester) async {
    final container = await pumpLogin(tester, FakeAuthRepository());

    await tester.tap(find.text('Masuk'));
    await tester.pump();

    expect(find.text('Email wajib diisi.'), findsOneWidget);
    expect(find.text('Password wajib diisi.'), findsOneWidget);
    expect(container.read(sessionControllerProvider).value, isNull);
  });

  testWidgets('signs in with valid credentials', (tester) async {
    final container = await pumpLogin(tester, FakeAuthRepository());

    await tester.enterText(
      find.byKey(const Key('login-email')),
      'rina@example.com',
    );
    await tester.enterText(
      find.byKey(const Key('login-password')),
      'rahasia123',
    );
    await tester.tap(find.text('Masuk'));
    await tester.pumpAndSettle();

    expect(container.read(sessionControllerProvider).value?.name, 'Rina');
  });

  testWidgets('shows the server message when credentials are wrong', (
    tester,
  ) async {
    await pumpLogin(
      tester,
      FakeAuthRepository(
        loginError: const ApiException(
          code: 'invalid_credentials',
          message: 'Email atau password salah.',
          statusCode: 401,
        ),
      ),
    );

    await tester.enterText(
      find.byKey(const Key('login-email')),
      'rina@example.com',
    );
    await tester.enterText(find.byKey(const Key('login-password')), 'salah123');
    await tester.tap(find.text('Masuk'));
    await tester.pumpAndSettle();

    expect(find.text('Email atau password salah.'), findsOneWidget);
  });
}
