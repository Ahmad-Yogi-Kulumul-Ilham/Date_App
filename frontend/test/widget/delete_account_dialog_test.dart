import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/network/api_exception.dart';
import 'package:leo_chat/core/theme/app_theme.dart';
import 'package:leo_chat/features/auth/data/api_auth_repository.dart';
import 'package:leo_chat/features/auth/presentation/session_controller.dart';
import 'package:leo_chat/features/profile/presentation/delete_account_dialog.dart';

import '../helpers/fakes.dart';

void main() {
  Future<ProviderContainer> openDialog(
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
          home: Builder(
            builder: (context) => Scaffold(
              body: TextButton(
                onPressed: () => showDeleteAccountDialog(context),
                child: const Text('buka'),
              ),
            ),
          ),
        ),
      ),
    );
    await container.read(sessionControllerProvider.future);
    await tester.tap(find.text('buka'));
    await tester.pumpAndSettle();
    return container;
  }

  final password = find.byKey(const Key('delete-account-password'));

  testWidgets('asks for the password before calling the API', (tester) async {
    final auth = FakeAuthRepository(restored: testProfile());
    await openDialog(tester, auth);

    await tester.tap(find.text('Hapus permanen'));
    await tester.pump();

    expect(find.text('Masukkan password untuk melanjutkan.'), findsOneWidget);
    expect(auth.deletedWithPassword, isNull);
  });

  testWidgets('shows the server message for a wrong password', (tester) async {
    final auth = FakeAuthRepository(restored: testProfile())
      ..deleteAccountError = const ApiException(
        code: 'wrong_password',
        message: 'Password salah.',
        statusCode: 422,
      );
    final container = await openDialog(tester, auth);

    await tester.enterText(password, 'salah');
    await tester.tap(find.text('Hapus permanen'));
    await tester.pumpAndSettle();

    expect(find.text('Password salah.'), findsOneWidget);
    expect(container.read(sessionControllerProvider).value, isNotNull);
  });

  testWidgets('deletes the account and closes', (tester) async {
    final auth = FakeAuthRepository(restored: testProfile());
    final container = await openDialog(tester, auth);

    await tester.enterText(password, 'rahasia123');
    await tester.tap(find.text('Hapus permanen'));
    await tester.pumpAndSettle();

    expect(auth.deletedWithPassword, 'rahasia123');
    expect(container.read(sessionControllerProvider).value, isNull);
    expect(find.text('Hapus akun?'), findsNothing);
  });
}
