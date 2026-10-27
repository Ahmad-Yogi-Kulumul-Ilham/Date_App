import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/location/location_service.dart';
import '../../../core/network/api_exception.dart';
import '../../../core/providers.dart';
import '../../../core/router/routes.dart';
import '../../../core/widgets/brand_logo.dart';
import '../../../core/widgets/state_views.dart';
import '../../matches/domain/match_summary.dart';
import '../domain/candidate.dart';
import 'discover_controller.dart';
import 'widgets/candidate_card.dart';
import 'widgets/match_dialog.dart';
import 'widgets/swipe_actions.dart';

typedef _OnSwipe = void Function(Candidate candidate, SwipeDirection direction);

class DiscoverPage extends ConsumerWidget {
  const DiscoverPage({super.key});

  // Runs with the page's context rather than the card's: swiping the last card
  // removes the card stack, but the match dialog must still be shown.
  Future<void> _swipe(
    BuildContext context,
    WidgetRef ref,
    Candidate candidate,
    SwipeDirection direction,
  ) async {
    try {
      final outcome = await ref
          .read(discoverControllerProvider.notifier)
          .swipe(candidate, direction);
      if (!outcome.matched || !context.mounted) return;

      final startChat = await showMatchDialog(context, candidate);
      if (startChat && context.mounted) {
        await context.push<void>(
          Routes.chat(outcome.matchId!),
          extra: MatchPartner(
            id: candidate.id,
            name: candidate.name,
            age: candidate.age,
            photoUrl: candidate.photoUrls.firstOrNull,
          ),
        );
      }
    } on ApiException catch (e) {
      if (context.mounted) showErrorSnackBar(context, e.message);
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final candidates = ref.watch(discoverControllerProvider);

    return Scaffold(
      appBar: AppBar(title: const BrandLogo()),
      body: SafeArea(
        child: switch (candidates) {
          AsyncData(:final value) when value.isEmpty => MessageView(
            icon: Icons.travel_explore_rounded,
            title: 'Belum ada orang baru',
            message:
                'Semua orang di sekitarmu sudah kamu lihat. Coba lagi nanti.',
            actionLabel: 'Muat ulang',
            onAction: () => ref.invalidate(discoverControllerProvider),
          ),
          AsyncData(:final value) => _CardStack(
            candidates: value,
            onSwipe: (c, d) => _swipe(context, ref, c, d),
          ),
          AsyncError(:final error) => _DiscoverError(error: error),
          _ => const LoadingView(),
        },
      ),
    );
  }
}

class _CardStack extends StatelessWidget {
  const _CardStack({required this.candidates, required this.onSwipe});

  final List<Candidate> candidates;
  final _OnSwipe onSwipe;

  @override
  Widget build(BuildContext context) {
    final top = candidates.first;
    final next = candidates.length > 1 ? candidates[1] : null;

    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 4, 16, 20),
      child: Column(
        children: [
          Expanded(
            child: Stack(
              fit: StackFit.expand,
              children: [
                if (next != null)
                  Transform.scale(
                    scale: 0.95,
                    child: CandidateCard(candidate: next),
                  ),
                Dismissible(
                  key: ValueKey(top.id),
                  onDismissed: (direction) => onSwipe(
                    top,
                    direction == DismissDirection.startToEnd
                        ? SwipeDirection.like
                        : SwipeDirection.pass,
                  ),
                  child: CandidateCard(candidate: top),
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),
          SwipeActions(
            onPass: () => onSwipe(top, SwipeDirection.pass),
            onLike: () => onSwipe(top, SwipeDirection.like),
          ),
        ],
      ),
    );
  }
}

class _DiscoverError extends ConsumerWidget {
  const _DiscoverError({required this.error});

  final Object error;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    void retry() => ref.invalidate(discoverControllerProvider);

    if (error case LocationException(:final reason, :final message)) {
      final permanent = reason == LocationFailure.permissionDeniedForever;
      return MessageView(
        icon: Icons.location_off_outlined,
        title: 'Lokasi dibutuhkan',
        message: message,
        actionLabel: permanent ? 'Buka pengaturan' : 'Coba lagi',
        onAction: permanent
            ? () => ref.read(locationServiceProvider).openSettings()
            : retry,
      );
    }

    return MessageView(
      icon: Icons.wifi_off_rounded,
      title: 'Gagal memuat',
      message: errorMessage(error),
      actionLabel: 'Coba lagi',
      onAction: retry,
    );
  }
}
