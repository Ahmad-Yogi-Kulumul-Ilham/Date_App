import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/providers.dart';
import '../domain/candidate.dart';
import '../domain/discover_repository.dart';

final discoverRepositoryProvider = Provider<DiscoverRepository>(
  (ref) => ApiDiscoverRepository(ref.watch(apiClientProvider)),
);

final class ApiDiscoverRepository implements DiscoverRepository {
  const ApiDiscoverRepository(this._api);

  final ApiClient _api;

  @override
  Future<List<Candidate>> fetchCandidates() async {
    final json = await _api.get<List<dynamic>>('/discover');
    return [
      for (final item in json) Candidate.fromJson(item as Map<String, dynamic>),
    ];
  }

  @override
  Future<SwipeOutcome> swipe(int userId, SwipeDirection direction) async {
    final json = await _api.post<Map<String, dynamic>>(
      '/swipes',
      body: {'user_id': userId, 'direction': direction.name},
    );
    return SwipeOutcome.fromJson(json);
  }
}
