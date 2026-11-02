import 'report_reason.dart';

abstract interface class SafetyRepository {
  /// Hides the two people from each other everywhere and ends their match.
  Future<void> block(int userId);

  /// Reports the user to moderators. The backend also blocks them.
  Future<void> report(int userId, ReportReason reason, {String details = ''});
}
