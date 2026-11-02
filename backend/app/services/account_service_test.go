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

func TestAccountService_Delete(t *testing.T) {
	ctx := context.Background()
	newOwner := func() *models.User {
		u := user(7, models.GenderMale, models.PreferenceFemale, time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC))
		u.Password = "hashed:rahasia123"
		u.Photos = []models.UserPhoto{{Path: "photos/7/a.jpg"}, {Path: "photos/7/b.jpg"}}
		return u
	}

	t.Run("deletes the user and their photo files", func(t *testing.T) {
		users := newFakeUsers(newOwner())
		files := &fakeFiles{}
		svc := NewAccountService(users, fakeHasher{}, files)

		require.NoError(t, svc.Delete(ctx, 7, "rahasia123"))

		_, err := users.FindByID(ctx, 7)
		assert.ErrorIs(t, err, apperr.ErrUserNotFound)
		assert.Equal(t, []string{"photos/7/a.jpg", "photos/7/b.jpg"}, files.deleted)
	})

	t.Run("a wrong password keeps the account and its files", func(t *testing.T) {
		users := newFakeUsers(newOwner())
		files := &fakeFiles{}
		svc := NewAccountService(users, fakeHasher{}, files)

		err := svc.Delete(ctx, 7, "salah")

		assert.ErrorIs(t, err, apperr.ErrWrongPassword)
		_, err = users.FindByID(ctx, 7)
		assert.NoError(t, err)
		assert.Empty(t, files.deleted)
	})
}
