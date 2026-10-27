package services

import (
	"context"
	"time"

	"leochat/app/geo"
	"leochat/app/models"
)

// The interfaces below are what the services need from the outside world.
// Production implementations live in app/repositories; tests use in-memory fakes.

type UserRepository interface {
	// FindByID returns apperr.ErrUserNotFound when no user exists.
	FindByID(ctx context.Context, id uint) (*models.User, error)
	// FindByEmail returns (nil, nil) when no user exists.
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Save(ctx context.Context, user *models.User) error
}

type PhotoRepository interface {
	ListByUser(ctx context.Context, userID uint) ([]models.UserPhoto, error)
	Create(ctx context.Context, photo *models.UserPhoto) error
	// Delete removes the photo only when it belongs to userID and returns the
	// deleted row, or apperr.ErrPhotoNotFound.
	Delete(ctx context.Context, userID, photoID uint) (*models.UserPhoto, error)
}

type CandidateQuery struct {
	ViewerID       uint
	Origin         geo.Point
	Box            geo.BoundingBox
	RadiusMeters   float64
	BornAfter      time.Time
	BornOnOrBefore time.Time
	// Genders the viewer wants to see.
	Genders []models.Gender
	// The viewer's own gender; candidates must be interested in it.
	ViewerGender models.Gender
	Limit        int
}

type Candidate struct {
	User           models.User
	DistanceMeters float64
}

type CandidateRepository interface {
	FindCandidates(ctx context.Context, q CandidateQuery) ([]Candidate, error)
}

type SwipeRepository interface {
	Exists(ctx context.Context, swiperID, swipedID uint) (bool, error)
	HasLiked(ctx context.Context, swiperID, swipedID uint) (bool, error)
	Create(ctx context.Context, swipe *models.Swipe) error
}

type MatchSummary struct {
	MatchID   uint
	MatchedAt time.Time
	// Partner has ID, Name, BirthDate and at most one photo (the primary) loaded.
	Partner     models.User
	LastMessage *MessagePreview
}

type MessagePreview struct {
	ID       uint
	SenderID uint
	Body     string
	SentAt   time.Time
	ReadAt   *time.Time
}

type MatchRepository interface {
	// FirstOrCreate returns the existing match for the pair or creates it.
	FirstOrCreate(ctx context.Context, userA, userB uint) (*models.Match, error)
	// FindByID returns apperr.ErrMatchNotFound when no match exists.
	FindByID(ctx context.Context, id uint) (*models.Match, error)
	ListForUser(ctx context.Context, userID uint) ([]MatchSummary, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	// ListByMatch returns up to limit messages older than beforeID (0 = newest),
	// newest first.
	ListByMatch(ctx context.Context, matchID, beforeID uint, limit int) ([]models.Message, error)
	MarkReadUpTo(ctx context.Context, matchID, readerID, upToID uint, at time.Time) error
}

type PasswordHasher interface {
	Make(value string) (string, error)
	Check(value, hashedValue string) bool
}

// FileUpload is an uploaded file that can persist itself under a directory and
// report the stored path. Goravel's filesystem.File satisfies it.
type FileUpload interface {
	Store(dir string) (string, error)
}

type FileRemover interface {
	Delete(paths ...string) error
}

// Clock lets tests control "now".
type Clock func() time.Time
