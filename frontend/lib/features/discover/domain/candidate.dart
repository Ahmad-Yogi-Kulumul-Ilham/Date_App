import '../../profile/domain/user_profile.dart';

/// Someone shown on the discover stack.
final class Candidate {
  const Candidate({
    required this.id,
    required this.name,
    required this.age,
    required this.gender,
    required this.bio,
    required this.distanceKm,
    required this.photoUrls,
  });

  factory Candidate.fromJson(Map<String, dynamic> json) => Candidate(
    id: json['id'] as int,
    name: json['name'] as String,
    age: json['age'] as int,
    gender: Gender.fromJson(json['gender'] as String),
    bio: json['bio'] as String? ?? '',
    distanceKm: json['distance_km'] as int,
    photoUrls: [
      for (final p in json['photos'] as List<dynamic>? ?? const [])
        (p as Map<String, dynamic>)['url'] as String,
    ],
  );

  final int id;
  final String name;
  final int age;
  final Gender gender;
  final String bio;
  final int distanceKm;
  final List<String> photoUrls;
}

enum SwipeDirection { like, pass }

final class SwipeOutcome {
  const SwipeOutcome({required this.matched, this.matchId});

  factory SwipeOutcome.fromJson(Map<String, dynamic> json) => SwipeOutcome(
    matched: json['matched'] as bool,
    matchId: json['match_id'] as int?,
  );

  final bool matched;
  final int? matchId;
}
