import 'package:flutter/material.dart';

import '../theme/app_colors.dart';

class BrandLogo extends StatelessWidget {
  const BrandLogo({super.key, this.size = 28});

  final double size;

  @override
  Widget build(BuildContext context) {
    return ShaderMask(
      shaderCallback: AppColors.brandGradient.createShader,
      blendMode: BlendMode.srcIn,
      child: Text(
        'leo',
        style: Theme.of(context).textTheme.headlineMedium?.copyWith(
          fontSize: size,
          fontWeight: FontWeight.w800,
          letterSpacing: -0.5,
        ),
      ),
    );
  }
}
