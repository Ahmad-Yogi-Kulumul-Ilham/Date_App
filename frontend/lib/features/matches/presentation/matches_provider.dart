import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../auth/presentation/current_user.dart';
import '../data/api_matches_repository.dart';
import '../domain/match_summary.dart';

final matchesProvider = FutureProvider.autoDispose<List<MatchSummary>>((ref) {
  ref.watch(currentUserIdProvider);
  return ref.watch(matchesRepositoryProvider).fetchMatches();
});
