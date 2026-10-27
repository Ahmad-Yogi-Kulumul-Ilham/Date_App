import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_exception.dart';
import '../../auth/presentation/validators.dart';
import '../domain/user_profile.dart';
import 'profile_controller.dart';

Future<void> showEditProfileSheet(BuildContext context, UserProfile profile) {
  return showModalBottomSheet<void>(
    context: context,
    isScrollControlled: true,
    showDragHandle: true,
    builder: (_) => _EditProfileSheet(profile: profile),
  );
}

class _EditProfileSheet extends ConsumerStatefulWidget {
  const _EditProfileSheet({required this.profile});

  final UserProfile profile;

  @override
  ConsumerState<_EditProfileSheet> createState() => _EditProfileSheetState();
}

class _EditProfileSheetState extends ConsumerState<_EditProfileSheet> {
  static const _maxBioLength = 500;

  final _formKey = GlobalKey<FormState>();
  late final _name = TextEditingController(text: widget.profile.name);
  late final _bio = TextEditingController(text: widget.profile.bio);
  late Preference _interestedIn = widget.profile.interestedIn;
  ApiException? _error;

  @override
  void dispose() {
    _name.dispose();
    _bio.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _error = null);
    try {
      await ref
          .read(profileControllerProvider.notifier)
          .updateDetails(
            name: _name.text,
            bio: _bio.text,
            interestedIn: _interestedIn,
          );
      if (mounted) Navigator.of(context).pop();
    } on ApiException catch (e) {
      if (mounted) setState(() => _error = e);
    }
  }

  @override
  Widget build(BuildContext context) {
    final busy = ref.watch(profileControllerProvider);
    final fields = _error?.fieldErrors ?? const {};

    return Padding(
      padding: EdgeInsets.fromLTRB(
        20,
        0,
        20,
        MediaQuery.viewInsetsOf(context).bottom + 24,
      ),
      child: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              'Edit profil',
              style: Theme.of(context).textTheme.titleLarge
                  ?.copyWith(fontWeight: FontWeight.w800),
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _name,
              decoration: const InputDecoration(labelText: 'Nama panggilan'),
              textCapitalization: TextCapitalization.words,
              validator: Validators.name,
              forceErrorText: fields['name'],
            ),
            const SizedBox(height: 14),
            TextFormField(
              controller: _bio,
              decoration: const InputDecoration(
                labelText: 'Bio',
                alignLabelWithHint: true,
              ),
              minLines: 3,
              maxLines: 5,
              maxLength: _maxBioLength,
              textCapitalization: TextCapitalization.sentences,
              forceErrorText: fields['bio'],
            ),
            const SizedBox(height: 8),
            Text(
              'Tertarik dengan',
              style: Theme.of(context).textTheme.labelLarge,
            ),
            const SizedBox(height: 8),
            SegmentedButton<Preference>(
              segments: [
                for (final p in Preference.values)
                  ButtonSegment(value: p, label: Text(p.label)),
              ],
              selected: {_interestedIn},
              onSelectionChanged: (s) =>
                  setState(() => _interestedIn = s.first),
            ),
            if (_error != null && fields.isEmpty) ...[
              const SizedBox(height: 16),
              Text(
                _error!.message,
                style: TextStyle(color: Theme.of(context).colorScheme.error),
              ),
            ],
            const SizedBox(height: 24),
            FilledButton(
              onPressed: busy ? null : _save,
              child: const Text('Simpan'),
            ),
          ],
        ),
      ),
    );
  }
}
