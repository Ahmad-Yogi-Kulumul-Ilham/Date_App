import 'package:dio/dio.dart';

/// An error the UI can show directly. [code] mirrors the backend's stable error
/// codes (for example `email_taken`), so callers can branch on it.
final class ApiException implements Exception {
  const ApiException({
    required this.code,
    required this.message,
    this.statusCode,
    this.fieldErrors = const {},
  });

  /// Converts any error thrown by Dio into an [ApiException], reading the
  /// backend's `{"error": {...}}` envelope when there is one.
  factory ApiException.fromDio(DioException error) {
    final response = error.response;
    final body = response?.data;
    if (body is Map<String, dynamic> && body['error'] is Map<String, dynamic>) {
      final err = body['error'] as Map<String, dynamic>;
      final fields = err['fields'];
      return ApiException(
        code: err['code'] as String? ?? 'unknown',
        message: err['message'] as String? ?? _fallbackMessage,
        statusCode: response?.statusCode,
        fieldErrors: fields is Map
            ? fields.map((k, v) => MapEntry(k.toString(), v.toString()))
            : const {},
      );
    }

    return switch (error.type) {
      DioExceptionType.connectionTimeout ||
      DioExceptionType.sendTimeout ||
      DioExceptionType.receiveTimeout ||
      DioExceptionType.connectionError => const ApiException(
        code: 'network',
        message: 'Tidak bisa terhubung ke server. Periksa koneksi internet.',
      ),
      _ => ApiException(
        code: 'unknown',
        message: _fallbackMessage,
        statusCode: response?.statusCode,
      ),
    };
  }

  static const _fallbackMessage = 'Terjadi kesalahan. Coba lagi nanti.';

  final String code;
  final String message;
  final int? statusCode;
  final Map<String, String> fieldErrors;

  bool get isUnauthorized => statusCode == 401;

  @override
  String toString() => 'ApiException($code, $statusCode): $message';
}

/// Message to show for any error, whether or not it came from the API.
String errorMessage(Object error) => error is ApiException
    ? error.message
    : 'Terjadi kesalahan. Coba lagi nanti.';
