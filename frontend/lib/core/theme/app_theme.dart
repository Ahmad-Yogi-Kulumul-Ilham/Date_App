import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import 'app_colors.dart';

abstract final class AppTheme {
  /// [useGoogleFonts] is off in widget tests, which have no network access.
  static ThemeData light({bool useGoogleFonts = true}) => _build(
    brightness: Brightness.light,
    primary: AppColors.plum,
    onPrimary: Colors.white,
    background: AppColors.lightBackground,
    surface: AppColors.lightSurface,
    text: AppColors.lightText,
    muted: AppColors.lightMuted,
    chip: AppColors.lightChip,
    chipText: AppColors.lightChipText,
    useGoogleFonts: useGoogleFonts,
  );

  static ThemeData dark({bool useGoogleFonts = true}) => _build(
    brightness: Brightness.dark,
    primary: AppColors.darkPrimary,
    onPrimary: AppColors.darkOnPrimary,
    background: AppColors.darkBackground,
    surface: AppColors.darkSurface,
    text: AppColors.darkText,
    muted: AppColors.darkMuted,
    chip: AppColors.darkChip,
    chipText: AppColors.darkChipText,
    useGoogleFonts: useGoogleFonts,
  );

  static ThemeData _build({
    required Brightness brightness,
    required Color primary,
    required Color onPrimary,
    required Color background,
    required Color surface,
    required Color text,
    required Color muted,
    required Color chip,
    required Color chipText,
    required bool useGoogleFonts,
  }) {
    final scheme =
        ColorScheme.fromSeed(
          seedColor: AppColors.plum,
          brightness: brightness,
        ).copyWith(
          primary: primary,
          onPrimary: onPrimary,
          secondary: AppColors.magenta,
          tertiary: AppColors.gold,
          surface: surface,
          onSurface: text,
          onSurfaceVariant: muted,
          surfaceContainerHighest: chip,
          error: AppColors.nope,
        );

    final base = ThemeData(useMaterial3: true, colorScheme: scheme);
    final textTheme =
        (useGoogleFonts
                ? GoogleFonts.plusJakartaSansTextTheme(base.textTheme)
                : base.textTheme)
            .apply(bodyColor: text, displayColor: text);

    return base.copyWith(
      scaffoldBackgroundColor: background,
      textTheme: textTheme,
      appBarTheme: AppBarTheme(
        backgroundColor: background,
        foregroundColor: text,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        titleTextStyle: textTheme.titleLarge?.copyWith(
          fontWeight: FontWeight.w800,
        ),
      ),
      cardTheme: CardThemeData(
        color: surface,
        elevation: 0,
        margin: EdgeInsets.zero,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
      ),
      chipTheme: base.chipTheme.copyWith(
        backgroundColor: chip,
        labelStyle: textTheme.labelMedium?.copyWith(
          color: chipText,
          fontWeight: FontWeight.w600,
        ),
        side: BorderSide.none,
        shape: const StadiumBorder(),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: onPrimary,
          minimumSize: const Size.fromHeight(52),
          shape: const StadiumBorder(),
          textStyle: textTheme.labelLarge?.copyWith(
            fontWeight: FontWeight.w700,
          ),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          minimumSize: const Size.fromHeight(52),
          shape: const StadiumBorder(),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: surface,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide.none,
        ),
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: surface,
        indicatorColor: chip,
      ),
    );
  }
}
