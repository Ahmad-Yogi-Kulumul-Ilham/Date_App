import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/utils/date_format.dart';

void main() {
  test('formatLongDate uses Indonesian month names', () {
    expect(formatLongDate(DateTime(2000, 5, 1)), '1 Mei 2000');
    expect(formatLongDate(DateTime(1998, 12, 31)), '31 Desember 1998');
  });

  group('formatListTimestamp', () {
    final now = DateTime(2026, 10, 26, 15);

    test('shows the time for today', () {
      expect(
        formatListTimestamp(DateTime(2026, 10, 26, 9, 5), now: now),
        '09.05',
      );
    });

    test('shows "Kemarin" for yesterday', () {
      expect(
        formatListTimestamp(DateTime(2026, 10, 25, 23), now: now),
        'Kemarin',
      );
    });

    test('shows day/month for older dates', () {
      expect(formatListTimestamp(DateTime(2026, 10, 3), now: now), '3/10');
    });
  });
}
