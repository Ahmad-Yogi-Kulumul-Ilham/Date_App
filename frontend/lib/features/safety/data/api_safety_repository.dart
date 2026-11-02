import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/providers.dart';
import '../domain/report_reason.dart';
import '../domain/safety_repository.dart';

final safetyRepositoryProvider = Provider<SafetyRepository>(
  (ref) => ApiSafetyRepository(ref.watch(apiClientProvider)),
);

final class ApiSafetyRepository implements SafetyRepository {
  const ApiSafetyRepository(this._api);

  final ApiClient _api;

  @override
  Future<void> block(int userId) =>
      _api.post<void>('/blocks', body: {'user_id': userId});

  @override
  Future<void> report(int userId, ReportReason reason, {String details = ''}) {
    final trimmed = details.trim();
    return _api.post<void>(
      '/reports',
      body: {
        'user_id': userId,
        'reason': reason.apiValue,
        if (trimmed.isNotEmpty) 'details': trimmed,
      },
    );
  }
}
