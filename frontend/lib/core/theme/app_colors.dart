import 'package:flutter/material.dart';

/// Plum Night palette.
abstract final class AppColors {
  static const plum = Color(0xFF7B2CBF);
  static const magenta = Color(0xFFE0479E);
  static const gold = Color(0xFFF9C74F);

  static const brandGradient = LinearGradient(
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
    colors: [plum, magenta],
  );

  // Swipe actions keep the same meaning in every theme.
  static const like = Color(0xFF22C55E);
  static const nope = Color(0xFFEF4444);
  static const superLike = Color(0xFF3B82F6);
  static const premium = Color(0xFFF59E0B);

  static const lightBackground = Color(0xFFFAF7FC);
  static const lightSurface = Color(0xFFFFFFFF);
  static const lightText = Color(0xFF1E1528);
  static const lightMuted = Color(0xFF6E6280);
  static const lightChip = Color(0xFFF1E6FA);
  static const lightChipText = Color(0xFF5A1F8C);

  static const darkBackground = Color(0xFF140F1A);
  static const darkSurface = Color(0xFF211A2B);
  static const darkText = Color(0xFFF1ECF7);
  static const darkMuted = Color(0xFFA99CB8);
  static const darkPrimary = Color(0xFFB57BEE);
  static const darkOnPrimary = Color(0xFF1A0F24);
  static const darkChip = Color(0xFF2E2340);
  static const darkChipText = Color(0xFFD9BFF5);
}
