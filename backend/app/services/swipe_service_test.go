package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"leochat/app/apperr"
	"leochat/app/models"
)

func newSwipeFixture() (*SwipeService, *fakeSwipes, *fakeMatches) {
	svc, swipes, matches, _ := newSwipeFixtureWithBlocks()
	return svc, swipes, matches
}

func newSwipeFixtureWithBlocks() (*SwipeService, *fakeSwipes, *fakeMatches, *fakeBlocks) {
	birth := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	users := newFakeUsers(
		user(1, models.GenderMale, models.PreferenceFemale, birth),
		user(2, models.GenderFemale, models.PreferenceMale, birth),
	)
	swipes := &fakeSwipes{}
	matches := &fakeMatches{}
	blocks := &fakeBlocks{}
	return NewSwipeService(users, swipes, matches, blocks), swipes, matches, blocks
}

func TestSwipeService_Swipe(t *testing.T) {
	ctx := context.Background()

	t.Run("one-sided like does not match", func(t *testing.T) {
		svc, _, matches := newSwipeFixture()

		result, err := svc.Swipe(ctx, 1, 2, models.SwipeLike)

		require.NoError(t, err)
		assert.False(t, result.Matched())
		assert.Empty(t, matches.rows)
	})

	t.Run("mutual like creates an ordered match", func(t *testing.T) {
		svc, swipes, _ := newSwipeFixture()
		swipes.rows = append(swipes.rows, models.Swipe{SwiperID: 2, SwipedID: 1, Direction: models.SwipeLike})

		result, err := svc.Swipe(ctx, 1, 2, models.SwipeLike)

		require.NoError(t, err)
		require.True(t, result.Matched())
		assert.Equal(t, uint(1), result.Match.UserOneID)
		assert.Equal(t, uint(2), result.Match.UserTwoID)
	})

	t.Run("pass never matches even if the other side liked", func(t *testing.T) {
		svc, swipes, _ := newSwipeFixture()
		swipes.rows = append(swipes.rows, models.Swipe{SwiperID: 2, SwipedID: 1, Direction: models.SwipeLike})

		result, err := svc.Swipe(ctx, 1, 2, models.SwipePass)

		require.NoError(t, err)
		assert.False(t, result.Matched())
	})

	t.Run("like after the other side passed does not match", func(t *testing.T) {
		svc, swipes, _ := newSwipeFixture()
		swipes.rows = append(swipes.rows, models.Swipe{SwiperID: 2, SwipedID: 1, Direction: models.SwipePass})

		result, err := svc.Swipe(ctx, 1, 2, models.SwipeLike)

		require.NoError(t, err)
		assert.False(t, result.Matched())
	})

	t.Run("rejects swiping yourself", func(t *testing.T) {
		svc, _, _ := newSwipeFixture()

		_, err := svc.Swipe(ctx, 1, 1, models.SwipeLike)

		assert.ErrorIs(t, err, apperr.ErrCannotSwipeSelf)
	})

	t.Run("rejects a second swipe on the same person", func(t *testing.T) {
		svc, _, _ := newSwipeFixture()
		_, err := svc.Swipe(ctx, 1, 2, models.SwipePass)
		require.NoError(t, err)

		_, err = svc.Swipe(ctx, 1, 2, models.SwipeLike)

		assert.ErrorIs(t, err, apperr.ErrAlreadySwiped)
	})

	t.Run("rejects unknown users", func(t *testing.T) {
		svc, _, _ := newSwipeFixture()

		_, err := svc.Swipe(ctx, 1, 99, models.SwipeLike)

		assert.ErrorIs(t, err, apperr.ErrUserNotFound)
	})

	t.Run("treats a blocked user as missing in both directions", func(t *testing.T) {
		svc, swipes, matches, blocks := newSwipeFixtureWithBlocks()
		swipes.rows = append(swipes.rows, models.Swipe{SwiperID: 1, SwipedID: 2, Direction: models.SwipeLike})
		blocks.rows = append(blocks.rows, models.Block{BlockerID: 1, BlockedID: 2})

		_, err := svc.Swipe(ctx, 2, 1, models.SwipeLike)

		assert.ErrorIs(t, err, apperr.ErrUserNotFound)
		assert.Len(t, swipes.rows, 1, "the blocked user's swipe is not recorded")
		assert.Empty(t, matches.rows)
	})
}
