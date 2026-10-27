import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

/// A person's photo, or their initial on the brand gradient when there is no
/// photo (or it fails to load).
class ProfilePhoto extends StatelessWidget {
  const ProfilePhoto({
    super.key,
    required this.name,
    this.url,
    this.initialSize = 72,
  });

  final String name;
  final String? url;
  final double initialSize;

  @override
  Widget build(BuildContext context) {
    final fallback = _InitialPlaceholder(name: name, fontSize: initialSize);
    final photoUrl = url;
    if (photoUrl == null) return fallback;

    return Image.network(
      photoUrl,
      fit: BoxFit.cover,
      errorBuilder: (_, _, _) => fallback,
      loadingBuilder: (context, child, progress) =>
          progress == null ? child : fallback,
    );
  }
}

/// Circular variant of [ProfilePhoto] for lists and headers.
class ProfileAvatar extends StatelessWidget {
  const ProfileAvatar({
    super.key,
    required this.name,
    this.url,
    this.radius = 24,
  });

  final String name;
  final String? url;
  final double radius;

  @override
  Widget build(BuildContext context) {
    return SizedBox.square(
      dimension: radius * 2,
      child: ClipOval(
        child: ProfilePhoto(name: name, url: url, initialSize: radius * 0.9),
      ),
    );
  }
}

class _InitialPlaceholder extends StatelessWidget {
  const _InitialPlaceholder({required this.name, required this.fontSize});

  final String name;
  final double fontSize;

  @override
  Widget build(BuildContext context) {
    final initial = name.trim().isEmpty ? '?' : name.trim()[0].toUpperCase();
    return DecoratedBox(
      decoration: const BoxDecoration(gradient: AppColors.brandGradient),
      child: Center(
        child: Text(
          initial,
          style: TextStyle(
            color: Colors.white.withValues(alpha: 0.85),
            fontSize: fontSize,
            fontWeight: FontWeight.w800,
          ),
        ),
      ),
    );
  }
}
