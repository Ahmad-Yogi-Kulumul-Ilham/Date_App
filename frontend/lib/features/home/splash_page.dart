import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/network/api_exception.dart';
import '../../core/widgets/brand_logo.dart';
import '../../core/widgets/state_views.dart';
import '../auth/presentation/session_controller.dart';

/// Shown while the saved session is restored, and if restoring it fails.
class SplashPage extends ConsumerWidget {
  const SplashPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final session = ref.watch(sessionControllerProvider);

    return Scaffold(
      body: SafeArea(
        child: session.hasError
            ? MessageView(
                icon: Icons.wifi_off_rounded,
                title: 'Tidak bisa terhubung',
                message: errorMessage(session.error!),
                actionLabel: 'Coba lagi',
                onAction: () => ref.invalidate(sessionControllerProvider),
              )
            : const Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    BrandLogo(size: 56),
                    SizedBox(height: 24),
                    SizedBox.square(
                      dimension: 24,
                      child: CircularProgressIndicator(strokeWidth: 2.5),
                    ),
                  ],
                ),
              ),
      ),
    );
  }
}
