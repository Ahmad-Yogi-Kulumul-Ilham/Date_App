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

func validRegistration() RegisterInput {
	return RegisterInput{
		Name:         "  Rina  ",
		Email:        "  Rina@Example.com ",
		Password:     "rahasia123",
		BirthDate:    time.Date(2000, time.May, 1, 0, 0, 0, 0, time.UTC),
		Gender:       models.GenderFemale,
		InterestedIn: models.PreferenceMale,
	}
}

func TestAuthService_Register(t *testing.T) {
	t.Run("normalizes and hashes", func(t *testing.T) {
		svc := NewAuthService(newFakeUsers(), fakeHasher{}, fixedClock)

		user, err := svc.Register(context.Background(), validRegistration())

		require.NoError(t, err)
		assert.Equal(t, "Rina", user.Name)
		assert.Equal(t, "rina@example.com", user.Email)
		assert.Equal(t, "hashed:rahasia123", user.Password)
		assert.NotZero(t, user.ID)
	})

	t.Run("rejects duplicate email regardless of case", func(t *testing.T) {
		existing := user(1, models.GenderMale, models.PreferenceFemale, time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC))
		existing.Email = "rina@example.com"
		svc := NewAuthService(newFakeUsers(existing), fakeHasher{}, fixedClock)

		_, err := svc.Register(context.Background(), validRegistration())

		assert.ErrorIs(t, err, apperr.ErrEmailTaken)
	})

	t.Run("rejects users under 18", func(t *testing.T) {
		in := validRegistration()
		// Turns 18 one day after fixedNow.
		in.BirthDate = time.Date(2008, time.October, 27, 0, 0, 0, 0, time.UTC)
		svc := NewAuthService(newFakeUsers(), fakeHasher{}, fixedClock)

		_, err := svc.Register(context.Background(), in)

		assert.ErrorIs(t, err, apperr.ErrUnderage)
	})

	t.Run("accepts users turning 18 today", func(t *testing.T) {
		in := validRegistration()
		in.BirthDate = time.Date(2008, time.October, 26, 0, 0, 0, 0, time.UTC)
		svc := NewAuthService(newFakeUsers(), fakeHasher{}, fixedClock)

		_, err := svc.Register(context.Background(), in)

		assert.NoError(t, err)
	})
}

func TestAuthService_Authenticate(t *testing.T) {
	stored := user(1, models.GenderFemale, models.PreferenceMale, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	stored.Email = "rina@example.com"
	stored.Password = "hashed:rahasia123"
	svc := NewAuthService(newFakeUsers(stored), fakeHasher{}, fixedClock)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{"correct credentials, email case-insensitive", "RINA@example.com", "rahasia123", nil},
		{"wrong password", "rina@example.com", "salah", apperr.ErrInvalidCredentials},
		{"unknown email gives the same error", "nobody@example.com", "rahasia123", apperr.ErrInvalidCredentials},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.Authenticate(context.Background(), tt.email, tt.password)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, stored.ID, got.ID)
		})
	}
}
