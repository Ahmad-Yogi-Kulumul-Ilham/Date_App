import 'candidate.dart';

abstract interface class DiscoverRepository {
  /// Nearby people not swiped yet, closest first.
  Future<List<Candidate>> fetchCandidates();

  Future<SwipeOutcome> swipe(int userId, SwipeDirection direction);
}
