import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/presentation/login_page.dart';
import '../../features/auth/presentation/register_page.dart';
import '../../features/auth/presentation/session_controller.dart';
import '../../features/chat/presentation/chat_page.dart';
import '../../features/discover/presentation/discover_page.dart';
import '../../features/home/home_shell.dart';
import '../../features/home/splash_page.dart';
import '../../features/matches/domain/match_summary.dart';
import '../../features/matches/presentation/matches_page.dart';
import '../../features/profile/domain/user_profile.dart';
import '../../features/profile/presentation/profile_page.dart';
import 'routes.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final session = ValueNotifier<AsyncValue<UserProfile?>>(const AsyncLoading());
  ref.listen(
    sessionControllerProvider,
    (_, next) => session.value = next,
    fireImmediately: true,
  );

  final router = GoRouter(
    initialLocation: Routes.discover,
    refreshListenable: session,
    redirect: (_, state) => redirectFor(session.value, state.matchedLocation),
    routes: [
      GoRoute(path: Routes.splash, builder: (_, _) => const SplashPage()),
      GoRoute(path: Routes.login, builder: (_, _) => const LoginPage()),
      GoRoute(path: Routes.register, builder: (_, _) => const RegisterPage()),
      GoRoute(
        path: Routes.chatPattern,
        builder: (_, state) => ChatPage(
          matchId: int.parse(state.pathParameters['matchId']!),
          partner: state.extra as MatchPartner?,
        ),
      ),
      StatefulShellRoute.indexedStack(
        builder: (_, _, shell) => HomeShell(shell: shell),
        branches: [
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: Routes.discover,
                builder: (_, _) => const DiscoverPage(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: Routes.matches,
                builder: (_, _) => const MatchesPage(),
              ),
            ],
          ),
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: Routes.profile,
                builder: (_, _) => const ProfilePage(),
              ),
            ],
          ),
        ],
      ),
    ],
  );

  ref.onDispose(() {
    router.dispose();
    session.dispose();
  });
  return router;
});

/// Where to send the user for [location] given the session state, or null to stay.
@visibleForTesting
String? redirectFor(AsyncValue<UserProfile?> session, String location) {
  final onAuthPage = location == Routes.login || location == Routes.register;

  // Still restoring the session, or restoring it failed (the splash offers a retry).
  if (!session.hasValue || session.hasError) {
    return location == Routes.splash ? null : Routes.splash;
  }

  final signedIn = session.value != null;
  if (!signedIn) return onAuthPage ? null : Routes.login;
  if (onAuthPage || location == Routes.splash) return Routes.discover;
  return null;
}
