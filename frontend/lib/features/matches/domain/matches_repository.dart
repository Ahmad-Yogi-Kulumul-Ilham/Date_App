import 'match_summary.dart';

abstract interface class MatchesRepository {
  /// Matches ordered by most recent activity.
  Future<List<MatchSummary>> fetchMatches();

  /// Ends the match for both people and deletes the conversation.
  Future<void> unmatch(int matchId);
}
