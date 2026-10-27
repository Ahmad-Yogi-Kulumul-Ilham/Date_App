import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/config/app_config.dart';
import '../data/api_chat_repository.dart';
import '../domain/chat_message.dart';
import '../domain/chat_repository.dart';

final chatControllerProvider = AsyncNotifierProvider.autoDispose
    .family<ChatController, List<ChatMessage>, int>(ChatController.new);

/// Messages of one conversation, newest first. While the chat is open it polls
/// for new messages; switching to WebSockets only needs to change this class.
class ChatController extends AsyncNotifier<List<ChatMessage>> {
  ChatController(this.matchId);

  final int matchId;

  /// Must match the backend's default page size.
  static const pageSize = 30;

  bool _hasOlder = true;
  bool _loadingOlder = false;

  bool get hasOlder => _hasOlder;

  ChatRepository get _repository => ref.read(chatRepositoryProvider);

  @override
  Future<List<ChatMessage>> build() async {
    final timer = Timer.periodic(AppConfig.chatPollInterval, (_) => _poll());
    ref.onDispose(timer.cancel);

    final page = await _repository.fetchMessages(matchId);
    _hasOlder = page.length >= pageSize;
    return page;
  }

  Future<void> send(String body) async {
    final message = await _repository.send(matchId, body.trim());
    if (ref.mounted) _mergeIn([message]);
  }

  Future<void> loadOlder() async {
    final oldest = state.value?.lastOrNull;
    if (!_hasOlder || _loadingOlder || oldest == null) return;

    _loadingOlder = true;
    try {
      final older = await _repository.fetchMessages(
        matchId,
        beforeId: oldest.id,
      );
      _hasOlder = older.length >= pageSize;
      if (ref.mounted) _mergeIn(older);
    } finally {
      _loadingOlder = false;
    }
  }

  Future<void> _poll() async {
    if (!state.hasValue) return;
    try {
      final latest = await _repository.fetchMessages(matchId);
      if (ref.mounted) _mergeIn(latest);
    } on Exception {
      // Polling failures are transient; the next tick tries again.
    }
  }

  void _mergeIn(List<ChatMessage> incoming) {
    final byId = {
      for (final m in state.value ?? const <ChatMessage>[]) m.id: m,
      for (final m in incoming) m.id: m,
    };
    state = AsyncData(
      byId.values.toList()..sort((a, b) => b.id.compareTo(a.id)),
    );
  }
}
