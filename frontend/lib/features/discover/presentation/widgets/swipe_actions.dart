import 'package:flutter/material.dart';

import '../../../../core/theme/app_colors.dart';

/// The pass / like buttons under the card stack.
class SwipeActions extends StatelessWidget {
  const SwipeActions({super.key, required this.onPass, required this.onLike});

  final VoidCallback? onPass;
  final VoidCallback? onLike;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _RoundButton(
          tooltip: 'Lewati',
          onPressed: onPass,
          size: 60,
          background: Theme.of(context).colorScheme.surface,
          child: const Icon(
            Icons.close_rounded,
            color: AppColors.nope,
            size: 30,
          ),
        ),
        const SizedBox(width: 28),
        _RoundButton(
          tooltip: 'Suka',
          onPressed: onLike,
          size: 72,
          gradient: AppColors.brandGradient,
          child: const Icon(
            Icons.favorite_rounded,
            color: Colors.white,
            size: 34,
          ),
        ),
      ],
    );
  }
}

class _RoundButton extends StatelessWidget {
  const _RoundButton({
    required this.tooltip,
    required this.onPressed,
    required this.size,
    required this.child,
    this.background,
    this.gradient,
  });

  final String tooltip;
  final VoidCallback? onPressed;
  final double size;
  final Widget child;
  final Color? background;
  final Gradient? gradient;

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: tooltip,
      child: Opacity(
        opacity: onPressed == null ? 0.5 : 1,
        child: Material(
          type: MaterialType.transparency,
          child: Ink(
            width: size,
            height: size,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: background,
              gradient: gradient,
              boxShadow: const [
                BoxShadow(
                  color: Colors.black26,
                  blurRadius: 12,
                  offset: Offset(0, 4),
                ),
              ],
            ),
            child: InkWell(
              customBorder: const CircleBorder(),
              onTap: onPressed,
              child: Center(child: child),
            ),
          ),
        ),
      ),
    );
  }
}
