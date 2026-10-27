import 'package:flutter/material.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/widgets/profile_photo.dart';
import '../../domain/candidate.dart';

/// Returns true when the user chose to start chatting.
Future<bool> showMatchDialog(BuildContext context, Candidate candidate) async {
  final result = await showDialog<bool>(
    context: context,
    builder: (context) => _MatchDialog(candidate: candidate),
  );
  return result ?? false;
}

class _MatchDialog extends StatelessWidget {
  const _MatchDialog({required this.candidate});

  final Candidate candidate;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    return Dialog(
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(28)),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(vertical: 28, horizontal: 24),
            decoration: const BoxDecoration(gradient: AppColors.brandGradient),
            child: Column(
              children: [
                ProfileAvatar(
                  name: candidate.name,
                  url: candidate.photoUrls.firstOrNull,
                  radius: 44,
                ),
                const SizedBox(height: 16),
                Text(
                  'Kalian cocok!',
                  style: textTheme.headlineSmall?.copyWith(
                    color: Colors.white,
                    fontWeight: FontWeight.w800,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  'Kamu dan ${candidate.name} saling suka.',
                  textAlign: TextAlign.center,
                  style: textTheme.bodyMedium?.copyWith(color: Colors.white),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              children: [
                FilledButton(
                  onPressed: () => Navigator.of(context).pop(true),
                  child: const Text('Kirim pesan'),
                ),
                const SizedBox(height: 8),
                TextButton(
                  onPressed: () => Navigator.of(context).pop(false),
                  child: const Text('Lanjut cari'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
