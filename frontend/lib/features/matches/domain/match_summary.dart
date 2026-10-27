final class MatchPartner {
  const MatchPartner({
    required this.id,
    required this.name,
    required this.age,
    this.photoUrl,
  });

  factory MatchPartner.fromJson(Map<String, dynamic> json) => MatchPartner(
    id: json['id'] as int,
    name: json['name'] as String,
    age: json['age'] as int,
    photoUrl: json['photo_url'] as String?,
  );

  final int id;
  final String name;
  final int age;
  final String? photoUrl;
}

final class LastMessage {
  const LastMessage({
    required this.body,
    required this.isMine,
    required this.sentAt,
    required this.isRead,
  });

  factory LastMessage.fromJson(Map<String, dynamic> json) => LastMessage(
    body: json['body'] as String,
    isMine: json['is_mine'] as bool,
    sentAt: DateTime.parse(json['sent_at'] as String),
    isRead: json['is_read'] as bool,
  );

  final String body;
  final bool isMine;
  final DateTime sentAt;
  final bool isRead;
}

final class MatchSummary {
  const MatchSummary({
    required this.id,
    required this.matchedAt,
    required this.partner,
    this.lastMessage,
  });

  factory MatchSummary.fromJson(Map<String, dynamic> json) => MatchSummary(
    id: json['id'] as int,
    matchedAt: DateTime.parse(json['matched_at'] as String),
    partner: MatchPartner.fromJson(json['partner'] as Map<String, dynamic>),
    lastMessage: json['last_message'] == null
        ? null
        : LastMessage.fromJson(json['last_message'] as Map<String, dynamic>),
  );

  final int id;
  final DateTime matchedAt;
  final MatchPartner partner;
  final LastMessage? lastMessage;

  /// A message from the partner the user has not opened yet.
  bool get hasUnread =>
      lastMessage != null && !lastMessage!.isMine && !lastMessage!.isRead;
}
