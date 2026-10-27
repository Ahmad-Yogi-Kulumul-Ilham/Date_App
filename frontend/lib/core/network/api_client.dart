import 'dart:async';

import 'package:dio/dio.dart';

import '../config/app_config.dart';
import '../storage/token_storage.dart';
import 'api_exception.dart';

/// Thin wrapper over Dio that knows the API's conventions: bearer auth, the
/// `{"data": ...}` envelope, error mapping and transparent token refresh.
class ApiClient {
  ApiClient({
    required this._tokens,
    String baseUrl = AppConfig.apiBaseUrl,
    Dio? dio,
  }) : _dio = dio ?? Dio() {
    _dio.options
      ..baseUrl = baseUrl
      ..connectTimeout = const Duration(seconds: 10)
      ..receiveTimeout = const Duration(seconds: 20)
      ..headers['Accept'] = 'application/json';
    _dio.interceptors.add(_AuthInterceptor(this));
  }

  final Dio _dio;
  final TokenStorage _tokens;
  final _sessionExpired = StreamController<void>.broadcast();

  /// Emits when an authenticated request fails and the token cannot be
  /// refreshed, meaning the user has to log in again.
  Stream<void> get sessionExpired => _sessionExpired.stream;

  Future<T> get<T>(String path, {Map<String, dynamic>? query}) =>
      _send(() => _dio.get<Object?>(path, queryParameters: query));

  Future<T> post<T>(String path, {Object? body}) =>
      _send(() => _dio.post<Object?>(path, data: body));

  Future<T> put<T>(String path, {Object? body}) =>
      _send(() => _dio.put<Object?>(path, data: body));

  Future<T> patch<T>(String path, {Object? body}) =>
      _send(() => _dio.patch<Object?>(path, data: body));

  Future<void> delete(String path) =>
      _send<Object?>(() => _dio.delete<Object?>(path));

  void dispose() => _sessionExpired.close();

  /// Returns the `data` field of the response body, or null for 204 responses.
  Future<T> _send<T>(Future<Response<Object?>> Function() request) async {
    try {
      final response = await request();
      final body = response.data;
      return (body is Map<String, dynamic> ? body['data'] : null) as T;
    } on DioException catch (e) {
      throw ApiException.fromDio(e);
    }
  }

  Future<String?> _refreshToken(String expiredToken) async {
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        '/auth/refresh',
        options: Options(
          headers: {'Authorization': 'Bearer $expiredToken'},
          extra: {_AuthInterceptor.skipAuthKey: true},
        ),
      );
      final data = response.data?['data'] as Map<String, dynamic>?;
      final token = data?['access_token'] as String?;
      if (token != null) await _tokens.write(token);
      return token;
    } on DioException {
      return null;
    }
  }

  Future<void> _expireSession() async {
    await _tokens.clear();
    if (!_sessionExpired.isClosed) _sessionExpired.add(null);
  }
}

/// Adds the bearer token and, on `token_expired`, refreshes once and replays
/// the request. Queued so that concurrent 401s trigger a single refresh.
class _AuthInterceptor extends QueuedInterceptor {
  _AuthInterceptor(this._client);

  static const skipAuthKey = 'skipAuth';
  static const _retriedKey = 'retried';

  final ApiClient _client;

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    if (options.extra[skipAuthKey] != true) {
      final token = await _client._tokens.read();
      if (token != null) options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    final options = err.requestOptions;
    final wasAuthenticated =
        options.extra[skipAuthKey] != true &&
        options.headers.containsKey('Authorization');
    if (err.response?.statusCode != 401 || !wasAuthenticated) {
      return handler.next(err);
    }

    final canRefresh =
        _errorCode(err.response) == 'token_expired' &&
        options.extra[_retriedKey] != true;
    final stale = await _client._tokens.read();
    final fresh = canRefresh && stale != null
        ? await _client._refreshToken(stale)
        : null;
    if (fresh == null) {
      await _client._expireSession();
      return handler.next(err);
    }

    try {
      options
        ..headers['Authorization'] = 'Bearer $fresh'
        ..extra[_retriedKey] = true;
      handler.resolve(await _client._dio.fetch<Object?>(options));
    } on DioException catch (e) {
      handler.next(e);
    }
  }

  static String? _errorCode(Response<Object?>? response) {
    final body = response?.data;
    if (body is Map<String, dynamic> && body['error'] is Map<String, dynamic>) {
      return (body['error'] as Map<String, dynamic>)['code'] as String?;
    }
    return null;
  }
}
