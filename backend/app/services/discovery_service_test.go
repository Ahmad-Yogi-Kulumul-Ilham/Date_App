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

func TestDiscoverFilter_Normalized(t *testing.T) {
	tests := []struct {
		name string
		in   DiscoverFilter
		want DiscoverFilter
	}{
		{"zero values use defaults", DiscoverFilter{}, DiscoverFilter{RadiusKm: 25, MinAge: 18, MaxAge: 60, Limit: 20}},
		{"values above limits are clamped", DiscoverFilter{RadiusKm: 999, MinAge: 30, MaxAge: 200, Limit: 500}, DiscoverFilter{RadiusKm: 150, MinAge: 30, MaxAge: 99, Limit: 50}},
		{"min age below 18 is raised", DiscoverFilter{MinAge: 12, MaxAge: 25}, DiscoverFilter{RadiusKm: 25, MinAge: 18, MaxAge: 25, Limit: 20}},
		{"max below min collapses to min", DiscoverFilter{MinAge: 30, MaxAge: 20}, DiscoverFilter{RadiusKm: 25, MinAge: 30, MaxAge: 30, Limit: 20}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.in.Normalized())
		})
	}
}

func TestDiscoveryService_Discover(t *testing.T) {
	ctx := context.Background()

	t.Run("requires the viewer's location", func(t *testing.T) {
		viewer := user(1, models.GenderMale, models.PreferenceFemale, time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC))
		svc := NewDiscoveryService(newFakeUsers(viewer), &fakeCandidates{}, fixedClock)

		_, err := svc.Discover(ctx, 1, DiscoverFilter{})

		assert.ErrorIs(t, err, apperr.ErrLocationRequired)
	})

	t.Run("builds the query from viewer preferences and filter", func(t *testing.T) {
		viewer := user(1, models.GenderFemale, models.PreferenceEveryone, time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC))
		viewer.Latitude, viewer.Longitude = ptr(-6.9), ptr(107.6)
		candidates := &fakeCandidates{}
		svc := NewDiscoveryService(newFakeUsers(viewer), candidates, fixedClock)

		_, err := svc.Discover(ctx, 1, DiscoverFilter{RadiusKm: 10, MinAge: 20, MaxAge: 30})
		require.NoError(t, err)

		q := candidates.got
		assert.Equal(t, uint(1), q.ViewerID)
		assert.Equal(t, 10_000.0, q.RadiusMeters)
		assert.ElementsMatch(t, []models.Gender{models.GenderMale, models.GenderFemale}, q.Genders)
		assert.Equal(t, models.GenderFemale, q.ViewerGender)
		// fixedNow is 2026-10-26: age 30 at most means born after 1995-10-26.
		assert.Equal(t, time.Date(1995, time.October, 26, 0, 0, 0, 0, time.UTC), q.BornAfter)
		// Age 20 at least means born on or before 2006-10-26.
		assert.Equal(t, time.Date(2006, time.October, 26, 0, 0, 0, 0, time.UTC), q.BornOnOrBefore)
		assert.Less(t, q.Box.MinLat, -6.9)
		assert.Greater(t, q.Box.MaxLng, 107.6)
	})
}
