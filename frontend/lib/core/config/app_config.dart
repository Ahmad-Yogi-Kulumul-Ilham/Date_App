/// Build-time configuration, set with `--dart-define`.
abstract final class AppConfig {
  /// Base URL of the backend API.
  ///
  /// Defaults to the Android emulator's alias for the host machine. For the iOS
  /// simulator run with
  /// `--dart-define=API_BASE_URL=http://localhost:3000/api/v1`.
  static const apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://10.0.2.2:3000/api/v1',
  );

  static const chatPollInterval = Duration(seconds: 4);
}
