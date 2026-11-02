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

func newSafetyFixture() (*SafetyService, *fakeBlocks, *fakeReports, *fakeMatches) {
	birth := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	users := newFakeUsers(
		user(1, models.GenderMale, models.PreferenceFemale, birth),
		user(2, models.GenderFemale, models.PreferenceMale, birth),
		user(3, models.GenderFemale, models.PreferenceMale, birth),
	)
	blocks := &fakeBlocks{}
	reports := &fakeReports{}
	matches := &fakeMatches{}
	return NewSafetyService(users, blocks, reports, matches), blocks, reports, matches
}

func TestSafetyService_Block(t *testing.T) {
	ctx := context.Background()

	t.Run("records the block and ends only that pair's match", func(t *testing.T) {
		svc, blocks, _, matches := newSafetyFixture()
		_, _ = matches.FirstOrCreate(ctx, 1, 2)
		_, _ = matches.FirstOrCreate(ctx, 1, 3)

		require.NoError(t, svc.Block(ctx, 1, 2))

		blocked, _ := blocks.ExistsBetween(ctx, 2, 1)
		assert.True(t, blocked, "a block works in both directions")
		require.Len(t, matches.rows, 1)
		assert.Equal(t, uint(3), matches.rows[0].UserTwoID)
	})

	t.Run("blocking the same user again is harmless", func(t *testing.T) {
		svc, blocks, _, _ := newSafetyFixture()

		require.NoError(t, svc.Block(ctx, 1, 2))
		require.NoError(t, svc.Block(ctx, 1, 2))

		assert.Len(t, blocks.rows, 1)
	})

	t.Run("rejects blocking yourself", func(t *testing.T) {
		svc, _, _, _ := newSafetyFixture()

		assert.ErrorIs(t, svc.Block(ctx, 1, 1), apperr.ErrCannotBlockSelf)
	})

	t.Run("rejects unknown users", func(t *testing.T) {
		svc, blocks, _, _ := newSafetyFixture()

		assert.ErrorIs(t, svc.Block(ctx, 1, 99), apperr.ErrUserNotFound)
		assert.Empty(t, blocks.rows)
	})
}

func TestSafetyService_Report(t *testing.T) {
	ctx := context.Background()

	t.Run("stores the report and blocks the reported user", func(t *testing.T) {
		svc, blocks, reports, matches := newSafetyFixture()
		_, _ = matches.FirstOrCreate(ctx, 1, 2)

		require.NoError(t, svc.Report(ctx, 1, 2, models.ReportHarassment, "  kasar di chat  "))

		require.Len(t, reports.rows, 1)
		assert.Equal(t, models.ReportHarassment, reports.rows[0].Reason)
		assert.Equal(t, "kasar di chat", reports.rows[0].Details)
		blocked, _ := blocks.ExistsBetween(ctx, 1, 2)
		assert.True(t, blocked)
		assert.Empty(t, matches.rows, "reporting also ends the match")
	})

	t.Run("rejects reporting yourself", func(t *testing.T) {
		svc, _, _, _ := newSafetyFixture()

		err := svc.Report(ctx, 1, 1, models.ReportSpam, "")

		assert.ErrorIs(t, err, apperr.ErrCannotReportSelf)
	})

	t.Run("rejects unknown users without storing anything", func(t *testing.T) {
		svc, blocks, reports, _ := newSafetyFixture()

		err := svc.Report(ctx, 1, 99, models.ReportSpam, "")

		assert.ErrorIs(t, err, apperr.ErrUserNotFound)
		assert.Empty(t, reports.rows)
		assert.Empty(t, blocks.rows)
	})
}
