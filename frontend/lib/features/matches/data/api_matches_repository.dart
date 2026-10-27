import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/providers.dart';
import '../domain/match_summary.dart';
import '../domain/matches_repository.dart';

final matchesRepositoryProvider = Provider<MatchesRepository>(
  (ref) => ApiMatchesRepository(ref.watch(apiClientProvider)),
);

final class ApiMatchesRepository implements MatchesRepository {
  const ApiMatchesRepository(this._api);

  final ApiClient _api;

  @override
  Future<List<MatchSummary>> fetchMatches() async {
    final json = await _api.get<List<dynamic>>('/matches');
    return [
      for (final item in json)
        MatchSummary.fromJson(item as Map<String, dynamic>),
    ];
  }
}
