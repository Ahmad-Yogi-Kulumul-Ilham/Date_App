import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/features/auth/domain/auth_repository.dart';
import 'package:leo_chat/features/discover/domain/candidate.dart';
import 'package:leo_chat/features/matches/domain/match_summary.dart';
import 'package:leo_chat/features/profile/domain/user_profile.dart';

void main() {
  test('UserProfile parses the /me payload', () {
    final profile = UserProfile.fromJson({
      'id': 1,
      'name': 'Rina',
      'email': 'rina@example.com',
      'birth_date': '2000-05-01',
      'age': 26,
      'gender': 'female',
      'interested_in': 'everyone',
      'bio': 'Halo',
      'has_location': true,
      'photos': [
        {'id': 7, 'url': 'http://x/7.jpg', 'position': 0},
      ],
    });

    expect(profile.birthDate, DateTime(2000, 5, 1));
    expect(profile.interestedIn, Preference.everyone);
    expect(profile.primaryPhotoUrl, 'http://x/7.jpg');
  });

  test('Candidate keeps only photo urls', () {
    final candidate = Candidate.fromJson({
      'id': 2,
      'name': 'Bima',
      'age': 27,
      'gender': 'male',
      'bio': '',
      'distance_km': 3,
      'photos': [
        {'id': 1, 'url': 'http://x/a.jpg', 'position': 0},
      ],
    });

    expect(candidate.photoUrls, ['http://x/a.jpg']);
  });

  test('MatchSummary flags unread messages from the partner only', () {
    Map<String, dynamic> json({required bool mine, required bool read}) => {
      'id': 1,
      'matched_at': '2026-10-26T09:00:00Z',
      'partner': {'id': 2, 'name': 'Bima', 'age': 27, 'photo_url': null},
      'last_message': {
        'body': 'Halo',
        'is_mine': mine,
        'sent_at': '2026-10-26T09:05:00Z',
        'is_read': read,
      },
    };

    expect(
      MatchSummary.fromJson(json(mine: false, read: false)).hasUnread,
      isTrue,
    );
    expect(
      MatchSummary.fromJson(json(mine: true, read: false)).hasUnread,
      isFalse,
    );
    expect(
      MatchSummary.fromJson(json(mine: false, read: true)).hasUnread,
      isFalse,
    );
  });

  test('RegisterData serializes dates as YYYY-MM-DD', () {
    final data = RegisterData(
      name: ' Rina ',
      email: 'rina@example.com',
      password: 'rahasia123',
      birthDate: DateTime(2000, 5, 1),
      gender: Gender.female,
      interestedIn: Preference.male,
    );

    expect(data.toJson()['birth_date'], '2000-05-01');
    expect(data.toJson()['name'], 'Rina');
  });
}
