import 'chat_message.dart';

abstract interface class ChatRepository {
  /// A page of messages, newest first. Pass [beforeId] to load older ones.
  Future<List<ChatMessage>> fetchMessages(int matchId, {int? beforeId});

  Future<ChatMessage> send(int matchId, String body);
}
