import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/providers.dart';
import '../domain/chat_message.dart';
import '../domain/chat_repository.dart';

final chatRepositoryProvider = Provider<ChatRepository>(
  (ref) => ApiChatRepository(ref.watch(apiClientProvider)),
);

final class ApiChatRepository implements ChatRepository {
  const ApiChatRepository(this._api);

  final ApiClient _api;

  @override
  Future<List<ChatMessage>> fetchMessages(int matchId, {int? beforeId}) async {
    final json = await _api.get<List<dynamic>>(
      '/matches/$matchId/messages',
      query: {'before_id': ?beforeId},
    );
    return [
      for (final item in json)
        ChatMessage.fromJson(item as Map<String, dynamic>),
    ];
  }

  @override
  Future<ChatMessage> send(int matchId, String body) async {
    final json = await _api.post<Map<String, dynamic>>(
      '/matches/$matchId/messages',
      body: {'body': body},
    );
    return ChatMessage.fromJson(json);
  }
}
