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

func TestProfileService_AddPhoto(t *testing.T) {
	ctx := context.Background()
	owner := user(7, models.GenderMale, models.PreferenceFemale, time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC))

	t.Run("stores under the user's folder and appends at the end", func(t *testing.T) {
		photos := &fakePhotos{rows: []models.UserPhoto{{UserID: 7, Position: 0}, {UserID: 7, Position: 3}}}
		svc := NewProfileService(newFakeUsers(owner), photos, &fakeFiles{}, fixedClock)
		upload := &fakeUpload{}

		photo, err := svc.AddPhoto(ctx, 7, upload)

		require.NoError(t, err)
		assert.Equal(t, "photos/7", upload.storedIn)
		assert.Equal(t, 4, photo.Position)
	})

	t.Run("enforces the per-user limit before storing anything", func(t *testing.T) {
		photos := &fakePhotos{}
		for i := range models.MaxPhotosPerUser {
			photos.rows = append(photos.rows, models.UserPhoto{UserID: 7, Position: i})
		}
		svc := NewProfileService(newFakeUsers(owner), photos, &fakeFiles{}, fixedClock)
		upload := &fakeUpload{}

		_, err := svc.AddPhoto(ctx, 7, upload)

		assert.ErrorIs(t, err, apperr.ErrPhotoLimitReached)
		assert.Empty(t, upload.storedIn)
	})
}

func TestProfileService_DeletePhoto(t *testing.T) {
	ctx := context.Background()
	photos := &fakePhotos{rows: []models.UserPhoto{{UserID: 7, Path: "photos/7/a.jpg"}}}
	photos.rows[0].ID = 1
	files := &fakeFiles{}
	svc := NewProfileService(newFakeUsers(), photos, files, fixedClock)

	assert.ErrorIs(t, svc.DeletePhoto(ctx, 8, 1), apperr.ErrPhotoNotFound, "another user's photo")

	require.NoError(t, svc.DeletePhoto(ctx, 7, 1))
	assert.Equal(t, []string{"photos/7/a.jpg"}, files.deleted)
}

func TestProfileService_Update(t *testing.T) {
	owner := user(7, models.GenderMale, models.PreferenceFemale, time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC))
	owner.Bio = "lama"
	svc := NewProfileService(newFakeUsers(owner), &fakePhotos{}, &fakeFiles{}, fixedClock)

	updated, err := svc.Update(context.Background(), 7, UpdateProfileInput{Name: ptr("  Bima "), InterestedIn: ptr(models.PreferenceEveryone)})

	require.NoError(t, err)
	assert.Equal(t, "Bima", updated.Name)
	assert.Equal(t, "lama", updated.Bio, "absent fields stay unchanged")
	assert.Equal(t, models.PreferenceEveryone, updated.InterestedIn)
}
