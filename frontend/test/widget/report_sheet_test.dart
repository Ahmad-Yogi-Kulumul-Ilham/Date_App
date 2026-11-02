import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/theme/app_theme.dart';
import 'package:leo_chat/features/safety/domain/report_reason.dart';
import 'package:leo_chat/features/safety/presentation/report_sheet.dart';

void main() {
  /// Opens the sheet and returns a getter for the value it resolves to.
  Future<ReportDraft? Function()> openSheet(WidgetTester tester) async {
    ReportDraft? result;
    await tester.pumpWidget(
      MaterialApp(
        theme: AppTheme.light(useGoogleFonts: false),
        home: Builder(
          builder: (context) => Scaffold(
            body: TextButton(
              onPressed: () async =>
                  result = await showReportSheet(context, name: 'Bima'),
              child: const Text('buka'),
            ),
          ),
        ),
      ),
    );
    await tester.tap(find.text('buka'));
    await tester.pumpAndSettle();
    return () => result;
  }

  final submit = find.widgetWithText(FilledButton, 'Kirim laporan');

  testWidgets('needs a reason before it can be sent', (tester) async {
    await openSheet(tester);

    expect(find.text('Laporkan Bima'), findsOneWidget);
    expect(tester.widget<FilledButton>(submit).onPressed, isNull);

    await tester.tap(find.text(ReportReason.harassment.label));
    await tester.pump();

    expect(tester.widget<FilledButton>(submit).onPressed, isNotNull);
  });

  testWidgets('returns the chosen reason and trimmed details', (tester) async {
    final result = await openSheet(tester);

    await tester.tap(find.text(ReportReason.spam.label));
    await tester.enterText(find.byType(TextField), '  jualan pulsa  ');
    await tester.ensureVisible(submit);
    await tester.tap(submit);
    await tester.pumpAndSettle();

    expect(result()?.reason, ReportReason.spam);
    expect(result()?.details, 'jualan pulsa');
    expect(find.text('Laporkan Bima'), findsNothing);
  });
}
