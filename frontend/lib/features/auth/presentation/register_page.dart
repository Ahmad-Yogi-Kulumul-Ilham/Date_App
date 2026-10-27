import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/network/api_exception.dart';
import '../../../core/router/routes.dart';
import '../../../core/utils/date_format.dart';
import '../../profile/domain/user_profile.dart';
import '../domain/auth_repository.dart';
import 'session_controller.dart';
import 'validators.dart';
import 'widgets/auth_layout.dart';

const _minimumAge = 18;

class RegisterPage extends ConsumerStatefulWidget {
  const RegisterPage({super.key});

  @override
  ConsumerState<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends ConsumerState<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  final _name = TextEditingController();
  final _email = TextEditingController();
  final _password = TextEditingController();
  final _birthDateText = TextEditingController();

  DateTime? _birthDate;
  Gender _gender = Gender.female;
  Preference _interestedIn = Preference.male;
  bool _submitting = false;
  ApiException? _error;

  @override
  void dispose() {
    _name.dispose();
    _email.dispose();
    _password.dispose();
    _birthDateText.dispose();
    super.dispose();
  }

  Future<void> _pickBirthDate() async {
    final now = DateTime.now();
    final latestAllowed = DateTime(now.year - _minimumAge, now.month, now.day);
    final picked = await showDatePicker(
      context: context,
      initialDate: _birthDate ?? DateTime(latestAllowed.year - 5),
      firstDate: DateTime(now.year - 99),
      lastDate: latestAllowed,
      helpText: 'Tanggal lahir',
    );
    if (picked == null) return;
    setState(() {
      _birthDate = picked;
      _birthDateText.text = formatLongDate(picked);
    });
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _submitting = true;
      _error = null;
    });
    try {
      await ref
          .read(sessionControllerProvider.notifier)
          .register(
            RegisterData(
              name: _name.text,
              email: _email.text,
              password: _password.text,
              birthDate: _birthDate!,
              gender: _gender,
              interestedIn: _interestedIn,
            ),
          );
    } on ApiException catch (e) {
      if (mounted) setState(() => _error = e);
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final fields = _error?.fieldErrors ?? const {};

    return AuthLayout(
      title: 'Buat akun',
      subtitle: 'Hanya butuh satu menit. Kamu harus berusia minimal 18 tahun.',
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextFormField(
              controller: _name,
              decoration: const InputDecoration(labelText: 'Nama panggilan'),
              textCapitalization: TextCapitalization.words,
              textInputAction: TextInputAction.next,
              validator: Validators.name,
              forceErrorText: fields['name'],
            ),
            const SizedBox(height: 14),
            TextFormField(
              controller: _email,
              decoration: const InputDecoration(labelText: 'Email'),
              keyboardType: TextInputType.emailAddress,
              textInputAction: TextInputAction.next,
              validator: Validators.email,
              forceErrorText: fields['email'],
            ),
            const SizedBox(height: 14),
            TextFormField(
              controller: _password,
              decoration: const InputDecoration(
                labelText: 'Password',
                helperText: 'Minimal 8 karakter',
              ),
              obscureText: true,
              validator: Validators.password,
              forceErrorText: fields['password'],
            ),
            const SizedBox(height: 14),
            TextFormField(
              controller: _birthDateText,
              decoration: const InputDecoration(
                labelText: 'Tanggal lahir',
                suffixIcon: Icon(Icons.calendar_today_outlined),
              ),
              readOnly: true,
              onTap: _pickBirthDate,
              validator: (_) =>
                  _birthDate == null ? 'Tanggal lahir wajib diisi.' : null,
              forceErrorText: fields['birth_date'],
            ),
            const SizedBox(height: 22),
            Text('Saya', style: theme.textTheme.labelLarge),
            const SizedBox(height: 8),
            SegmentedButton<Gender>(
              segments: [
                for (final g in Gender.values)
                  ButtonSegment(value: g, label: Text(g.label)),
              ],
              selected: {_gender},
              onSelectionChanged: (s) => setState(() => _gender = s.first),
            ),
            const SizedBox(height: 18),
            Text('Tertarik dengan', style: theme.textTheme.labelLarge),
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
            const SizedBox(height: 24),
            if (_error != null && fields.isEmpty) ...[
              FormErrorBanner(message: _error!.message),
              const SizedBox(height: 16),
            ],
            FilledButton(
              onPressed: _submitting ? null : _submit,
              child: _submitting
                  ? const SizedBox.square(
                      dimension: 22,
                      child: CircularProgressIndicator(strokeWidth: 2.5),
                    )
                  : const Text('Daftar'),
            ),
            const SizedBox(height: 12),
            TextButton(
              onPressed: _submitting ? null : () => context.go(Routes.login),
              child: const Text('Sudah punya akun? Masuk'),
            ),
          ],
        ),
      ),
    );
  }
}
