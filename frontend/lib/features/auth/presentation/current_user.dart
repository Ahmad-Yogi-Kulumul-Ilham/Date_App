import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'session_controller.dart';

/// The signed-in user's id. Providers holding per-user data watch this so they
/// reset automatically when someone else logs in.
final currentUserIdProvider = Provider<int?>(
  (ref) => ref.watch(sessionControllerProvider.select((s) => s.value?.id)),
);
