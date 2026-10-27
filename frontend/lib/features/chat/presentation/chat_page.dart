import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_exception.dart';
import '../../../core/utils/date_format.dart';
import '../../../core/widgets/profile_photo.dart';
import '../../../core/widgets/state_views.dart';
import '../../matches/domain/match_summary.dart';
import '../domain/chat_message.dart';
import 'chat_controller.dart';

class ChatPage extends ConsumerWidget {
  const ChatPage({super.key, required this.matchId, this.partner});

  final int matchId;

  /// Shown in the header; null when the page is opened without context
  /// (for example from a deep link).
  final MatchPartner? partner;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final messages = ref.watch(chatControllerProvider(matchId));
    final partner = this.partner;

    return Scaffold(
      appBar: AppBar(
        titleSpacing: 0,
        title: partner == null
            ? const Text('Chat')
            : Row(
                children: [
                  ProfileAvatar(
                    name: partner.name,
                    url: partner.photoUrl,
                    radius: 18,
                  ),
                  const SizedBox(width: 10),
                  Flexible(
                    child: Text(partner.name, overflow: TextOverflow.ellipsis),
                  ),
                ],
              ),
      ),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: switch (messages) {
                AsyncData(:final value) when value.isEmpty => MessageView(
                  icon: Icons.waving_hand_outlined,
                  title: 'Mulai obrolan',
                  message: partner == null
                      ? 'Kirim pesan pertama.'
                      : 'Kirim pesan pertama ke ${partner.name}.',
                ),
                AsyncData(:final value) => _MessageList(
                  matchId: matchId,
                  messages: value,
                ),
                AsyncError(:final error) => MessageView(
                  icon: Icons.wifi_off_rounded,
                  title: 'Gagal memuat pesan',
                  message: errorMessage(error),
                  actionLabel: 'Coba lagi',
                  onAction: () =>
                      ref.invalidate(chatControllerProvider(matchId)),
                ),
                _ => const LoadingView(),
              },
            ),
            _Composer(matchId: matchId),
          ],
        ),
      ),
    );
  }
}

class _MessageList extends ConsumerWidget {
  const _MessageList({required this.matchId, required this.messages});

  final int matchId;
  final List<ChatMessage> messages;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return NotificationListener<ScrollNotification>(
      onNotification: (n) {
        // The list is reversed, so maxScrollExtent is the oldest message.
        if (n.metrics.extentAfter < 200) {
          ref.read(chatControllerProvider(matchId).notifier).loadOlder();
        }
        return false;
      },
      child: ListView.builder(
        reverse: true,
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        itemCount: messages.length,
        itemBuilder: (context, i) => _Bubble(message: messages[i]),
      ),
    );
  }
}

class _Bubble extends StatelessWidget {
  const _Bubble({required this.message});

  final ChatMessage message;

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    final mine = message.isMine;
    final background = mine ? scheme.primary : scheme.surface;
    final foreground = mine ? scheme.onPrimary : scheme.onSurface;

    return Align(
      alignment: mine ? Alignment.centerRight : Alignment.centerLeft,
      child: ConstrainedBox(
        constraints: BoxConstraints(
          maxWidth: MediaQuery.sizeOf(context).width * 0.78,
        ),
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 3),
          padding: const EdgeInsets.fromLTRB(14, 9, 14, 7),
          decoration: BoxDecoration(
            color: background,
            borderRadius: BorderRadius.only(
              topLeft: const Radius.circular(18),
              topRight: const Radius.circular(18),
              bottomLeft: Radius.circular(mine ? 18 : 4),
              bottomRight: Radius.circular(mine ? 4 : 18),
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                message.body,
                style: TextStyle(color: foreground, height: 1.35),
              ),
              const SizedBox(height: 2),
              Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(
                    formatClock(message.sentAt),
                    style: TextStyle(
                      color: foreground.withValues(alpha: 0.7),
                      fontSize: 11,
                    ),
                  ),
                  if (mine) ...[
                    const SizedBox(width: 4),
                    Icon(
                      message.readAt == null
                          ? Icons.done_rounded
                          : Icons.done_all_rounded,
                      size: 14,
                      color: foreground.withValues(alpha: 0.7),
                    ),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _Composer extends ConsumerStatefulWidget {
  const _Composer({required this.matchId});

  final int matchId;

  @override
  ConsumerState<_Composer> createState() => _ComposerState();
}

class _ComposerState extends ConsumerState<_Composer> {
  static const _maxLength = 2000;

  final _controller = TextEditingController();
  bool _sending = false;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _send() async {
    final body = _controller.text.trim();
    if (body.isEmpty || _sending) return;

    setState(() => _sending = true);
    try {
      await ref
          .read(chatControllerProvider(widget.matchId).notifier)
          .send(body);
      _controller.clear();
    } on ApiException catch (e) {
      if (mounted) showErrorSnackBar(context, e.message);
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final scheme = Theme.of(context).colorScheme;
    return Padding(
      padding: const EdgeInsets.fromLTRB(12, 6, 8, 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Expanded(
            child: TextField(
              controller: _controller,
              minLines: 1,
              maxLines: 5,
              maxLength: _maxLength,
              textCapitalization: TextCapitalization.sentences,
              decoration: const InputDecoration(
                hintText: 'Tulis pesan...',
                counterText: '',
                contentPadding: EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 12,
                ),
              ),
              onSubmitted: (_) => _send(),
            ),
          ),
          const SizedBox(width: 6),
          IconButton.filled(
            tooltip: 'Kirim',
            onPressed: _sending ? null : _send,
            style: IconButton.styleFrom(backgroundColor: scheme.primary),
            icon: Icon(Icons.send_rounded, color: scheme.onPrimary),
          ),
        ],
      ),
    );
  }
}
