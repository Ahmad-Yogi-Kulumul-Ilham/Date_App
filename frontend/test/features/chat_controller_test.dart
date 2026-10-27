import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:leo_chat/features/chat/data/api_chat_repository.dart';
import 'package:leo_chat/features/chat/presentation/chat_controller.dart';

import '../helpers/fakes.dart';

void main() {
  late FakeChatRepository chat;
  late ProviderContainer container;

  setUp(() {
    // 40 messages, newest (id 40) first: more than one page.
    chat = FakeChatRepository([
      for (var id = 40; id >= 1; id--) testMessage(id),
    ]);
    container = ProviderContainer(
      overrides: [chatRepositoryProvider.overrideWithValue(chat)],
      retry: (_, _) => null,
    );
    addTearDown(container.dispose);
    container.listen(chatControllerProvider(1), (_, _) {});
  });

  List<int> ids() => [
    for (final m in container.read(chatControllerProvider(1)).value!) m.id,
  ];

  test('loads the newest page first', () async {
    final page = await container.read(chatControllerProvider(1).future);

    expect(page.first.id, 40);
    expect(page, hasLength(ChatController.pageSize));
    expect(container.read(chatControllerProvider(1).notifier).hasOlder, isTrue);
  });

  test('loadOlder appends older messages without duplicates', () async {
    await container.read(chatControllerProvider(1).future);
    final controller = container.read(chatControllerProvider(1).notifier);

    await controller.loadOlder();

    expect(ids(), [for (var id = 40; id >= 1; id--) id]);
    expect(controller.hasOlder, isFalse);
  });

  test('sent messages appear at the top', () async {
    await container.read(chatControllerProvider(1).future);

    await container.read(chatControllerProvider(1).notifier).send('  Halo!  ');

    final newest = container.read(chatControllerProvider(1)).value!.first;
    expect(newest.body, 'Halo!');
    expect(newest.isMine, isTrue);
  });
}
