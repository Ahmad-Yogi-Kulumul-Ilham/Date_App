final class ChatMessage {
  const ChatMessage({
    required this.id,
    required this.matchId,
    required this.senderId,
    required this.body,
    required this.isMine,
    required this.sentAt,
    this.readAt,
  });

  factory ChatMessage.fromJson(Map<String, dynamic> json) => ChatMessage(
    id: json['id'] as int,
    matchId: json['match_id'] as int,
    senderId: json['sender_id'] as int,
    body: json['body'] as String,
    isMine: json['is_mine'] as bool,
    sentAt: DateTime.parse(json['sent_at'] as String),
    readAt: json['read_at'] == null
        ? null
        : DateTime.parse(json['read_at'] as String),
  );

  final int id;
  final int matchId;
  final int senderId;
  final String body;
  final bool isMine;
  final DateTime sentAt;
  final DateTime? readAt;
}
