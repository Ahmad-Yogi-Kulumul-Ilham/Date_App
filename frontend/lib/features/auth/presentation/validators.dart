/// Client-side checks that mirror the backend's rules, so most mistakes are
/// caught before a request is sent.
abstract final class Validators {
  static final _email = RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$');

  static String? email(String? value) {
    final v = value?.trim() ?? '';
    if (v.isEmpty) return 'Email wajib diisi.';
    if (!_email.hasMatch(v)) return 'Format email belum benar.';
    return null;
  }

  static String? password(String? value) {
    final v = value ?? '';
    if (v.isEmpty) return 'Password wajib diisi.';
    if (v.length < 8) return 'Password minimal 8 karakter.';
    if (v.length > 72) return 'Password maksimal 72 karakter.';
    return null;
  }

  static String? name(String? value) {
    final v = value?.trim() ?? '';
    if (v.length < 2) return 'Nama minimal 2 karakter.';
    if (v.length > 50) return 'Nama maksimal 50 karakter.';
    return null;
  }

  static String? required(String? value, String label) =>
      (value?.trim().isEmpty ?? true) ? '$label wajib diisi.' : null;
}
