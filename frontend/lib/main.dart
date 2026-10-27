import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'app.dart';

void main() {
  runApp(
    // Errors are shown to the user with a retry button instead of being
    // retried silently in the background.
    ProviderScope(retry: (_, _) => null, child: const LeoChatApp()),
  );
}
