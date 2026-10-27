import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/core/network/api_exception.dart';

DioException _dioError({int? status, Object? body, DioExceptionType? type}) {
  final options = RequestOptions(path: '/x');
  return DioException(
    requestOptions: options,
    type: type ?? DioExceptionType.badResponse,
    response: status == null
        ? null
        : Response<Object?>(
            requestOptions: options,
            statusCode: status,
            data: body,
          ),
  );
}

void main() {
  test('reads the backend error envelope', () {
    final e = ApiException.fromDio(
      _dioError(
        status: 422,
        body: {
          'error': {
            'code': 'validation_failed',
            'message': 'Data yang dikirim belum valid.',
            'fields': {'email': 'Email sudah dipakai.'},
          },
        },
      ),
    );

    expect(e.code, 'validation_failed');
    expect(e.statusCode, 422);
    expect(e.fieldErrors, {'email': 'Email sudah dipakai.'});
  });

  test('maps connection problems to a network error', () {
    final e = ApiException.fromDio(
      _dioError(type: DioExceptionType.connectionError),
    );

    expect(e.code, 'network');
    expect(e.message, contains('koneksi'));
  });

  test('falls back to a generic error for unexpected bodies', () {
    final e = ApiException.fromDio(_dioError(status: 502, body: '<html>'));

    expect(e.code, 'unknown');
    expect(e.statusCode, 502);
  });
}
