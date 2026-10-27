import 'package:flutter/material.dart';

import '../../../../core/widgets/profile_photo.dart';
import '../../domain/candidate.dart';

class CandidateCard extends StatelessWidget {
  const CandidateCard({super.key, required this.candidate});

  final Candidate candidate;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    const onPhoto = Colors.white;

    return Card(
      clipBehavior: Clip.antiAlias,
      elevation: 6,
      shadowColor: Colors.black38,
      child: Stack(
        fit: StackFit.expand,
        children: [
          ProfilePhoto(
            name: candidate.name,
            url: candidate.photoUrls.firstOrNull,
            initialSize: 140,
          ),
          const DecoratedBox(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                stops: [0.45, 1],
                colors: [Colors.transparent, Color(0xCC000000)],
              ),
            ),
          ),
          Positioned(
            left: 20,
            right: 20,
            bottom: 20,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text.rich(
                  TextSpan(
                    children: [
                      TextSpan(
                        text: candidate.name,
                        style: const TextStyle(fontWeight: FontWeight.w800),
                      ),
                      TextSpan(text: ', ${candidate.age}'),
                    ],
                  ),
                  style: textTheme.headlineSmall?.copyWith(color: onPhoto),
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    const Icon(Icons.place_outlined, size: 16, color: onPhoto),
                    const SizedBox(width: 4),
                    Text(
                      '${candidate.distanceKm} km dari kamu',
                      style: textTheme.bodySmall?.copyWith(color: onPhoto),
                    ),
                  ],
                ),
                if (candidate.bio.isNotEmpty) ...[
                  const SizedBox(height: 10),
                  Text(
                    candidate.bio,
                    maxLines: 3,
                    overflow: TextOverflow.ellipsis,
                    style: textTheme.bodyMedium?.copyWith(color: onPhoto),
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}
