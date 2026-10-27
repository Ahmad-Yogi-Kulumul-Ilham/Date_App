abstract final class Routes {
  static const splash = '/splash';
  static const login = '/login';
  static const register = '/register';
  static const discover = '/discover';
  static const matches = '/matches';
  static const profile = '/profile';
  static const chatPattern = '/chat/:matchId';

  static String chat(int matchId) => '/chat/$matchId';
}
