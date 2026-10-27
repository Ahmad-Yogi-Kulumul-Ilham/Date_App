package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"leochat/app/apperr"
	"leochat/app/models"
)

const MinimumAge = 18

type AuthService struct {
	users  UserRepository
	hasher PasswordHasher
	now    Clock
}

func NewAuthService(users UserRepository, hasher PasswordHasher, now Clock) *AuthService {
	return &AuthService{users: users, hasher: hasher, now: now}
}

type RegisterInput struct {
	Name         string
	Email        string
	Password     string
	BirthDate    time.Time
	Gender       models.Gender
	InterestedIn models.Preference
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*models.User, error) {
	if models.YearsBetween(in.BirthDate, s.now()) < MinimumAge {
		return nil, apperr.ErrUnderage
	}

	email := normalizeEmail(in.Email)
	existing, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperr.ErrEmailTaken
	}

	hashed, err := s.hasher.Make(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Name:         strings.TrimSpace(in.Name),
		Email:        email,
		Password:     hashed,
		BirthDate:    in.BirthDate,
		Gender:       in.Gender,
		InterestedIn: in.InterestedIn,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Authenticate verifies credentials. Unknown email and wrong password return the
// same error so the API does not reveal which emails are registered.
func (s *AuthService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, err
	}
	if user == nil || !s.hasher.Check(password, user.Password) {
		return nil, apperr.ErrInvalidCredentials
	}

	return user, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
