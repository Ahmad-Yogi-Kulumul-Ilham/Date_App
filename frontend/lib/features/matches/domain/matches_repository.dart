import 'match_summary.dart';

abstract interface class MatchesRepository {
  /// Matches ordered by most recent activity.
  Future<List<MatchSummary>> fetchMatches();
}
