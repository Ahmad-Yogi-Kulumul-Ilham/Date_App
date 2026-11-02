import 'dart:async';

import 'package:leo_chat/core/location/location_service.dart';
import 'package:leo_chat/core/network/api_exception.dart';
import 'package:leo_chat/features/auth/domain/auth_repository.dart';
import 'package:leo_chat/features/chat/domain/chat_message.dart';
import 'package:leo_chat/features/chat/domain/chat_repository.dart';
import 'package:leo_chat/features/discover/domain/candidate.dart';
import 'package:leo_chat/features/discover/domain/discover_repository.dart';
import 'package:leo_chat/features/matches/domain/match_summary.dart';
import 'package:leo_chat/features/matches/domain/matches_repository.dart';
import 'package:leo_chat/features/profile/domain/profile_repository.dart';
import 'package:leo_chat/features/profile/domain/user_profile.dart';
import 'package:leo_chat/features/safety/domain/report_reason.dart';
import 'package:leo_chat/features/safety/domain/safety_repository.dart';

UserProfile testProfile({int id = 1, String name = 'Rina'}) => UserProfile(
  id: id,
  name: name,
  email: '${name.toLowerCase()}@example.com',
  birthDate: DateTime(2000, 5, 1),
  age: 26,
  gender: Gender.female,
  interestedIn: Preference.male,
  bio: '',
  hasLocation: true,
  photos: const [],
);

Candidate testCandidate(int id, {String name = 'Bima'}) => Candidate(
  id: id,
  name: name,
  age: 27,
  gender: Gender.male,
  bio: 'Suka futsal.',
  distanceKm: 3,
  photoUrls: const [],
);

ChatMessage testMessage(int id, {bool mine = false}) => ChatMessage(
  id: id,
  matchId: 1,
  senderId: mine ? 1 : 2,
  body: 'pesan $id',
  isMine: mine,
  sentAt: DateTime.utc(2026, 10, 26, 9, id),
);

class FakeAuthRepository implements AuthRepository {
  FakeAuthRepository({this.restored, this.loginError});

  UserProfile? restored;
  ApiException? loginError;
  int logoutCalls = 0;
  final _expired = StreamController<void>.broadcast();

  void expireSession() => _expired.add(null);

  @override
  Stream<void> get sessionExpired => _expired.stream;

  @override
  Future<UserProfile?> restoreSession() async => restored;

  @override
  Future<UserProfile> login({
    required String email,
    required String password,
  }) async {
    if (loginError != null) throw loginError!;
    return testProfile();
  }

  @override
  Future<UserProfile> register(RegisterData data) async =>
      testProfile(name: data.name);

  @override
  Future<void> logout() async => logoutCalls++;
}

class FakeLocationService implements LocationService {
  FakeLocationService({this.failure});

  LocationFailure? failure;

  @override
  Future<Coordinates> currentPosition() async {
    if (failure != null) throw LocationException(failure!);
    return const Coordinates(-6.9, 107.6);
  }

  @override
  Future<void> openSettings() async {}
}

class FakeProfileRepository implements ProfileRepository {
  Coordinates? lastLocation;

  @override
  Future<void> updateLocation(Coordinates coordinates) async =>
      lastLocation = coordinates;

  @override
  Future<UserProfile> updateProfile({
    String? name,
    String? bio,
    Preference? interestedIn,
  }) async => testProfile(name: name ?? 'Rina');

  @override
  Future<Photo> uploadPhoto(String filePath) async =>
      const Photo(id: 1, url: 'http://x/1.jpg', position: 0);

  @override
  Future<void> deletePhoto(int photoId) async {}
}

class FakeDiscoverRepository implements DiscoverRepository {
  FakeDiscoverRepository(this.candidates);

  List<Candidate> candidates;
  final Set<int> likesBack = {};
  ApiException? swipeError;
  final swipes = <(int, SwipeDirection)>[];

  @override
  Future<List<Candidate>> fetchCandidates() async => candidates;

  @override
  Future<SwipeOutcome> swipe(int userId, SwipeDirection direction) async {
    if (swipeError != null) throw swipeError!;
    swipes.add((userId, direction));
    final matched =
        direction == SwipeDirection.like && likesBack.contains(userId);
    return SwipeOutcome(matched: matched, matchId: matched ? 99 : null);
  }
}

class FakeMatchesRepository implements MatchesRepository {
  int fetchCalls = 0;
  final unmatched = <int>[];

  @override
  Future<List<MatchSummary>> fetchMatches() async {
    fetchCalls++;
    return const [];
  }

  @override
  Future<void> unmatch(int matchId) async => unmatched.add(matchId);
}

class FakeSafetyRepository implements SafetyRepository {
  final blocked = <int>[];
  final reports = <(int, ReportReason, String)>[];
  ApiException? error;

  @override
  Future<void> block(int userId) async {
    if (error != null) throw error!;
    blocked.add(userId);
  }

  @override
  Future<void> report(
    int userId,
    ReportReason reason, {
    String details = '',
  }) async {
    if (error != null) throw error!;
    reports.add((userId, reason, details));
  }
}

class FakeChatRepository implements ChatRepository {
  FakeChatRepository(this.messages);

  /// All messages in the conversation, newest first.
  List<ChatMessage> messages;
  int nextId = 1000;

  @override
  Future<List<ChatMessage>> fetchMessages(int matchId, {int? beforeId}) async {
    final older = beforeId == null
        ? messages
        : messages.where((m) => m.id < beforeId).toList();
    return older.take(30).toList();
  }

  @override
  Future<ChatMessage> send(int matchId, String body) async {
    final message = ChatMessage(
      id: nextId++,
      matchId: matchId,
      senderId: 1,
      body: body,
      isMine: true,
      sentAt: DateTime.utc(2026, 10, 26, 10),
    );
    messages = [message, ...messages];
    return message;
  }
}
