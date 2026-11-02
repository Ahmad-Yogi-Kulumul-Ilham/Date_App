import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

import '../../../core/network/api_exception.dart';
import '../../../core/widgets/confirm_dialog.dart';
import '../../../core/widgets/profile_photo.dart';
import '../../../core/widgets/state_views.dart';
import '../../auth/presentation/session_controller.dart';
import '../domain/user_profile.dart';
import 'delete_account_dialog.dart';
import 'edit_profile_sheet.dart';
import 'profile_controller.dart';

class ProfilePage extends ConsumerWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final profile = ref.watch(sessionControllerProvider).value;
    if (profile == null) return const Scaffold(body: LoadingView());

    final theme = Theme.of(context);
    final busy = ref.watch(profileControllerProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Profil'),
        actions: [
          IconButton(
            tooltip: 'Edit profil',
            icon: const Icon(Icons.edit_outlined),
            onPressed: () => showEditProfileSheet(context, profile),
          ),
        ],
        bottom: busy
            ? const PreferredSize(
                preferredSize: Size.fromHeight(2),
                child: LinearProgressIndicator(minHeight: 2),
              )
            : null,
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 32),
        children: [
          Row(
            children: [
              ProfileAvatar(
                name: profile.name,
                url: profile.primaryPhotoUrl,
                radius: 36,
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${profile.name}, ${profile.age}',
                      style: theme.textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      profile.email,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          _SectionTitle(
            'Foto (${profile.photos.length}/${UserProfile.maxPhotos})',
          ),
          const SizedBox(height: 10),
          _PhotoGrid(profile: profile, enabled: !busy),
          const SizedBox(height: 6),
          Text(
            'Ketuk foto untuk menghapus. Foto pertama jadi foto utama.',
            style: theme.textTheme.bodySmall?.copyWith(
              color: theme.colorScheme.onSurfaceVariant,
            ),
          ),
          const SizedBox(height: 24),
          const _SectionTitle('Tentang saya'),
          const SizedBox(height: 8),
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Text(
                profile.bio.isEmpty
                    ? 'Belum ada bio. Ceritakan sedikit tentang dirimu.'
                    : profile.bio,
                style: theme.textTheme.bodyMedium,
              ),
            ),
          ),
          const SizedBox(height: 16),
          Card(
            child: Column(
              children: [
                ListTile(
                  leading: const Icon(Icons.person_outline),
                  title: const Text('Saya'),
                  trailing: Text(profile.gender.label),
                ),
                ListTile(
                  leading: const Icon(Icons.favorite_outline),
                  title: const Text('Tertarik dengan'),
                  trailing: Text(profile.interestedIn.label),
                ),
              ],
            ),
          ),
          const SizedBox(height: 32),
          OutlinedButton.icon(
            onPressed: () => _confirmLogout(context, ref),
            icon: const Icon(Icons.logout),
            label: const Text('Keluar'),
          ),
          const SizedBox(height: 8),
          TextButton(
            onPressed: () => showDeleteAccountDialog(context),
            style: TextButton.styleFrom(
              foregroundColor: theme.colorScheme.error,
            ),
            child: const Text('Hapus akun'),
          ),
        ],
      ),
    );
  }

  Future<void> _confirmLogout(BuildContext context, WidgetRef ref) async {
    final confirmed = await showConfirmDialog(
      context,
      title: 'Keluar dari akun?',
      confirmLabel: 'Keluar',
    );
    if (confirmed) {
      await ref.read(sessionControllerProvider.notifier).logout();
    }
  }
}

class _SectionTitle extends StatelessWidget {
  const _SectionTitle(this.text);

  final String text;

  @override
  Widget build(BuildContext context) => Text(
    text,
    style: Theme.of(context).textTheme.titleSmall
        ?.copyWith(fontWeight: FontWeight.w700),
  );
}

class _PhotoGrid extends ConsumerWidget {
  const _PhotoGrid({required this.profile, required this.enabled});

  final UserProfile profile;
  final bool enabled;

  Future<void> _add(BuildContext context, WidgetRef ref) async {
    final file = await ImagePicker().pickImage(
      source: ImageSource.gallery,
      maxWidth: 1440,
      imageQuality: 85,
    );
    if (file == null) return;

    try {
      await ref.read(profileControllerProvider.notifier).addPhoto(file.path);
    } on ApiException catch (e) {
      if (context.mounted) showErrorSnackBar(context, e.message);
    }
  }

  Future<void> _delete(BuildContext context, WidgetRef ref, Photo photo) async {
    final confirmed = await showConfirmDialog(
      context,
      title: 'Hapus foto ini?',
      confirmLabel: 'Hapus',
    );
    if (!confirmed) return;

    try {
      await ref.read(profileControllerProvider.notifier).deletePhoto(photo);
    } on ApiException catch (e) {
      if (context.mounted) showErrorSnackBar(context, e.message);
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final scheme = Theme.of(context).colorScheme;
    final canAdd = profile.photos.length < UserProfile.maxPhotos;

    return GridView.count(
      crossAxisCount: 3,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      mainAxisSpacing: 8,
      crossAxisSpacing: 8,
      childAspectRatio: 3 / 4,
      children: [
        for (final photo in profile.photos)
          InkWell(
            borderRadius: BorderRadius.circular(14),
            onTap: enabled ? () => _delete(context, ref, photo) : null,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(14),
              child: ProfilePhoto(
                name: profile.name,
                url: photo.url,
                initialSize: 32,
              ),
            ),
          ),
        if (canAdd)
          Material(
            color: scheme.surfaceContainerHighest,
            borderRadius: BorderRadius.circular(14),
            child: InkWell(
              borderRadius: BorderRadius.circular(14),
              onTap: enabled ? () => _add(context, ref) : null,
              child: Center(
                child: Icon(Icons.add_a_photo_outlined, color: scheme.primary),
              ),
            ),
          ),
      ],
    );
  }
}
