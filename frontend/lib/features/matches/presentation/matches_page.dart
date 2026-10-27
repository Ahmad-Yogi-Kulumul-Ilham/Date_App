import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/network/api_exception.dart';
import '../../../core/router/routes.dart';
import '../../../core/utils/date_format.dart';
import '../../../core/widgets/profile_photo.dart';
import '../../../core/widgets/state_views.dart';
import '../domain/match_summary.dart';
import 'matches_provider.dart';

class MatchesPage extends ConsumerWidget {
  const MatchesPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final matches = ref.watch(matchesProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Match')),
      body: switch (matches) {
        AsyncData(:final value) when value.isEmpty => const MessageView(
          icon: Icons.favorite_border_rounded,
          title: 'Belum ada match',
          message: 'Saat kamu dan seseorang saling suka, kalian bisa mulai ngobrol di sini.',
        ),
        AsyncData(:final value) => RefreshIndicator(
          onRefresh: () => ref.refresh(matchesProvider.future),
          child: ListView.separated(
            padding: const EdgeInsets.symmetric(vertical: 8),
            itemCount: value.length,
            separatorBuilder: (_, _) => const Divider(height: 1, indent: 84),
            itemBuilder: (context, i) => _MatchTile(
              match: value[i],
              onTap: () async {
                await context.push(
                  Routes.chat(value[i].id),
                  extra: value[i].partner,
                );
                ref.invalidate(matchesProvider);
              },
            ),
          ),
        ),
        AsyncError(:final error) => MessageView(
          icon: Icons.wifi_off_rounded,
          title: 'Gagal memuat',
          message: errorMessage(error),
          actionLabel: 'Coba lagi',
          onAction: () => ref.invalidate(matchesProvider),
        ),
        _ => const LoadingView(),
      },
    );
  }
}

class _MatchTile extends StatelessWidget {
  const _MatchTile({required this.match, required this.onTap});

  final MatchSummary match;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final last = match.lastMessage;
    final preview = switch (last) {
      null => 'Kalian baru match. Sapa duluan!',
      LastMessage(isMine: true) => 'Kamu: ${last.body}',
      _ => last.body,
    };

    return ListTile(
      onTap: onTap,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
      leading: ProfileAvatar(
        name: match.partner.name,
        url: match.partner.photoUrl,
        radius: 28,
      ),
      title: Text(
        '${match.partner.name}, ${match.partner.age}',
        style: theme.textTheme.titleMedium?.copyWith(
          fontWeight: FontWeight.w700,
        ),
      ),
      subtitle: Text(
        preview,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
        style: TextStyle(
          fontWeight: match.hasUnread ? FontWeight.w700 : FontWeight.normal,
          color: match.hasUnread
              ? theme.colorScheme.onSurface
              : theme.colorScheme.onSurfaceVariant,
        ),
      ),
      trailing: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Text(
            formatListTimestamp(last?.sentAt ?? match.matchedAt),
            style: theme.textTheme.labelSmall?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
          if (match.hasUnread) ...[
            const SizedBox(height: 6),
            Container(
              width: 10,
              height: 10,
              decoration: BoxDecoration(
                color: theme.colorScheme.secondary,
                shape: BoxShape.circle,
              ),
            ),
          ],
        ],
      ),
    );
  }
}
