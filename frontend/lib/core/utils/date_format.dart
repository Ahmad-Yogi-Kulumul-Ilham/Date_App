const _months = [
  'Januari',
  'Februari',
  'Maret',
  'April',
  'Mei',
  'Juni',
  'Juli',
  'Agustus',
  'September',
  'Oktober',
  'November',
  'Desember',
];

String _twoDigits(int n) => n.toString().padLeft(2, '0');

/// "19.42", the Indonesian clock format.
String formatClock(DateTime time) {
  final t = time.toLocal();
  return '${_twoDigits(t.hour)}.${_twoDigits(t.minute)}';
}

/// "1 Mei 2000".
String formatLongDate(DateTime date) =>
    '${date.day} ${_months[date.month - 1]} ${date.year}';

/// Compact timestamp for lists: the time for today, "Kemarin", otherwise "12/10".
String formatListTimestamp(DateTime time, {DateTime? now}) {
  final t = time.toLocal();
  final today = (now ?? DateTime.now()).toLocal();
  final dayDiff = DateTime(
    today.year,
    today.month,
    today.day,
  ).difference(DateTime(t.year, t.month, t.day)).inDays;

  return switch (dayDiff) {
    0 => formatClock(t),
    1 => 'Kemarin',
    _ => '${t.day}/${t.month}',
  };
}
