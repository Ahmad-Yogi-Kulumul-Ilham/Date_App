package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"leochat/app/apperr"
	"leochat/app/models"
)

func newChatFixture() (*ChatService, *fakeMessages) {
	matches := &fakeMatches{}
	_, _ = matches.FirstOrCreate(context.Background(), 1, 2)
	messages := &fakeMessages{}
	return NewChatService(matches, messages, fixedClock), messages
}

func TestChatService_Send(t *testing.T) {
	ctx := context.Background()

	t.Run("participant can send a trimmed message", func(t *testing.T) {
		svc, _ := newChatFixture()

		msg, err := svc.Send(ctx, 2, 1, "  Halo!  ")

		require.NoError(t, err)
		assert.Equal(t, "Halo!", msg.Body)
		assert.Equal(t, uint(2), msg.SenderID)
	})

	t.Run("outsider gets not found", func(t *testing.T) {
		svc, _ := newChatFixture()

		_, err := svc.Send(ctx, 3, 1, "Halo")

		assert.ErrorIs(t, err, apperr.ErrMatchNotFound)
	})

	t.Run("whitespace-only message is rejected", func(t *testing.T) {
		svc, _ := newChatFixture()

		_, err := svc.Send(ctx, 1, 1, "   ")

		assert.ErrorIs(t, err, apperr.ErrMessageEmpty)
	})
}

func TestChatService_Messages(t *testing.T) {
	ctx := context.Background()

	t.Run("newest page marks messages read up to the newest id", func(t *testing.T) {
		svc, messages := newChatFixture()
		for _, body := range []string{"a", "b", "c"} {
			_, err := svc.Send(ctx, 2, 1, body)
			require.NoError(t, err)
		}

		page, err := svc.Messages(ctx, 1, 1, 0, 2)

		require.NoError(t, err)
		require.Len(t, page, 2)
		assert.Equal(t, "c", page[0].Body)
		assert.Equal(t, uint(3), messages.markedUpTo)
	})

	t.Run("older pages do not touch read state", func(t *testing.T) {
		svc, messages := newChatFixture()
		for _, body := range []string{"a", "b", "c"} {
			_, err := svc.Send(ctx, 2, 1, body)
			require.NoError(t, err)
		}

		page, err := svc.Messages(ctx, 1, 1, 3, 10)

		require.NoError(t, err)
		assert.Equal(t, []string{"b", "a"}, bodies(page))
		assert.Zero(t, messages.markedUpTo)
	})
}

func bodies(messages []models.Message) []string {
	out := make([]string, len(messages))
	for i, m := range messages {
		out[i] = m.Body
	}
	return out
}
