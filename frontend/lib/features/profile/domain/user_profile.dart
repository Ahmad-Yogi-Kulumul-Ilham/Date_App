enum Gender {
  male('Pria'),
  female('Wanita');

  const Gender(this.label);

  final String label;

  static Gender fromJson(String value) => Gender.values.firstWhere(
    (g) => g.name == value,
    orElse: () => Gender.male,
  );
}

enum Preference {
  male('Pria'),
  female('Wanita'),
  everyone('Semua');

  const Preference(this.label);

  final String label;

  static Preference fromJson(String value) => Preference.values.firstWhere(
    (p) => p.name == value,
    orElse: () => Preference.everyone,
  );
}

final class Photo {
  const Photo({required this.id, required this.url, required this.position});

  factory Photo.fromJson(Map<String, dynamic> json) => Photo(
    id: json['id'] as int,
    url: json['url'] as String,
    position: json['position'] as int? ?? 0,
  );

  final int id;
  final String url;
  final int position;
}

/// The signed-in user's own profile.
final class UserProfile {
  const UserProfile({
    required this.id,
    required this.name,
    required this.email,
    required this.birthDate,
    required this.age,
    required this.gender,
    required this.interestedIn,
    required this.bio,
    required this.hasLocation,
    required this.photos,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) => UserProfile(
    id: json['id'] as int,
    name: json['name'] as String,
    email: json['email'] as String,
    birthDate: DateTime.parse(json['birth_date'] as String),
    age: json['age'] as int,
    gender: Gender.fromJson(json['gender'] as String),
    interestedIn: Preference.fromJson(json['interested_in'] as String),
    bio: json['bio'] as String? ?? '',
    hasLocation: json['has_location'] as bool? ?? false,
    photos: [
      for (final p in json['photos'] as List<dynamic>? ?? const [])
        Photo.fromJson(p as Map<String, dynamic>),
    ],
  );

  final int id;
  final String name;
  final String email;
  final DateTime birthDate;
  final int age;
  final Gender gender;
  final Preference interestedIn;
  final String bio;
  final bool hasLocation;
  final List<Photo> photos;

  static const maxPhotos = 6;

  String? get primaryPhotoUrl => photos.isEmpty ? null : photos.first.url;

  UserProfile copyWith({List<Photo>? photos}) => UserProfile(
    id: id,
    name: name,
    email: email,
    birthDate: birthDate,
    age: age,
    gender: gender,
    interestedIn: interestedIn,
    bio: bio,
    hasLocation: hasLocation,
    photos: photos ?? this.photos,
  );
}
