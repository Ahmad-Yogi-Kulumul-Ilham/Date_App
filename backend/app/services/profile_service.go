package services

import (
	"context"
	"fmt"
	"strings"

	"leochat/app/apperr"
	"leochat/app/geo"
	"leochat/app/models"
)

type ProfileService struct {
	users  UserRepository
	photos PhotoRepository
	files  FileRemover
	now    Clock
}

func NewProfileService(users UserRepository, photos PhotoRepository, files FileRemover, now Clock) *ProfileService {
	return &ProfileService{users: users, photos: photos, files: files, now: now}
}

func (s *ProfileService) Get(ctx context.Context, userID uint) (*models.User, error) {
	return s.users.FindByID(ctx, userID)
}

// UpdateProfileInput holds optional changes; nil fields are left untouched.
type UpdateProfileInput struct {
	Name         *string
	Bio          *string
	InterestedIn *models.Preference
}

func (s *ProfileService) Update(ctx context.Context, userID uint, in UpdateProfileInput) (*models.User, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		user.Name = strings.TrimSpace(*in.Name)
	}
	if in.Bio != nil {
		user.Bio = strings.TrimSpace(*in.Bio)
	}
	if in.InterestedIn != nil {
		user.InterestedIn = *in.InterestedIn
	}

	if err := s.users.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ProfileService) UpdateLocation(ctx context.Context, userID uint, point geo.Point) error {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	now := s.now()
	user.Latitude = &point.Lat
	user.Longitude = &point.Lng
	user.LastActiveAt = &now

	return s.users.Save(ctx, user)
}

func (s *ProfileService) AddPhoto(ctx context.Context, userID uint, upload FileUpload) (*models.UserPhoto, error) {
	existing, err := s.photos.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(existing) >= models.MaxPhotosPerUser {
		return nil, apperr.ErrPhotoLimitReached
	}

	path, err := upload.Store(fmt.Sprintf("photos/%d", userID))
	if err != nil {
		return nil, fmt.Errorf("store photo: %w", err)
	}

	photo := &models.UserPhoto{UserID: userID, Path: path, Position: nextPosition(existing)}
	if err := s.photos.Create(ctx, photo); err != nil {
		_ = s.files.Delete(path)
		return nil, err
	}

	return photo, nil
}

func (s *ProfileService) DeletePhoto(ctx context.Context, userID, photoID uint) error {
	photo, err := s.photos.Delete(ctx, userID, photoID)
	if err != nil {
		return err
	}

	// The row is already gone; a leftover file is harmless and must not fail the request.
	_ = s.files.Delete(photo.Path)
	return nil
}

func nextPosition(photos []models.UserPhoto) int {
	next := 0
	for _, p := range photos {
		if p.Position >= next {
			next = p.Position + 1
		}
	}
	return next
}
