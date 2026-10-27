package services

import (
	"context"
	"time"

	"leochat/app/apperr"
	"leochat/app/models"
)

// In-memory implementations of the ports, just enough to exercise service logic.

var fixedNow = time.Date(2026, time.October, 26, 9, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedNow }

type fakeUsers struct {
	byID   map[uint]*models.User
	nextID uint
}

func newFakeUsers(users ...*models.User) *fakeUsers {
	f := &fakeUsers{byID: map[uint]*models.User{}}
	for _, u := range users {
		f.byID[u.ID] = u
		f.nextID = max(f.nextID, u.ID)
	}
	return f
}

func (f *fakeUsers) FindByID(_ context.Context, id uint) (*models.User, error) {
	if u, ok := f.byID[id]; ok {
		return u, nil
	}
	return nil, apperr.ErrUserNotFound
}

func (f *fakeUsers) FindByEmail(_ context.Context, email string) (*models.User, error) {
	for _, u := range f.byID {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (f *fakeUsers) Create(_ context.Context, u *models.User) error {
	f.nextID++
	u.ID = f.nextID
	f.byID[u.ID] = u
	return nil
}

func (f *fakeUsers) Save(_ context.Context, u *models.User) error {
	f.byID[u.ID] = u
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Make(v string) (string, error) { return "hashed:" + v, nil }
func (fakeHasher) Check(v, h string) bool        { return h == "hashed:"+v }

type fakeSwipes struct {
	rows []models.Swipe
}

func (f *fakeSwipes) Exists(_ context.Context, swiperID, swipedID uint) (bool, error) {
	for _, s := range f.rows {
		if s.SwiperID == swiperID && s.SwipedID == swipedID {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeSwipes) HasLiked(_ context.Context, swiperID, swipedID uint) (bool, error) {
	for _, s := range f.rows {
		if s.SwiperID == swiperID && s.SwipedID == swipedID && s.Direction == models.SwipeLike {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeSwipes) Create(_ context.Context, s *models.Swipe) error {
	f.rows = append(f.rows, *s)
	return nil
}

type fakeMatches struct {
	rows []models.Match
}

func (f *fakeMatches) FirstOrCreate(_ context.Context, a, b uint) (*models.Match, error) {
	one, two := models.OrderedPair(a, b)
	for i := range f.rows {
		if f.rows[i].UserOneID == one && f.rows[i].UserTwoID == two {
			return &f.rows[i], nil
		}
	}
	m := models.Match{UserOneID: one, UserTwoID: two}
	m.ID = uint(len(f.rows) + 1)
	f.rows = append(f.rows, m)
	return &f.rows[len(f.rows)-1], nil
}

func (f *fakeMatches) FindByID(_ context.Context, id uint) (*models.Match, error) {
	for i := range f.rows {
		if f.rows[i].ID == id {
			return &f.rows[i], nil
		}
	}
	return nil, apperr.ErrMatchNotFound
}

func (f *fakeMatches) ListForUser(context.Context, uint) ([]MatchSummary, error) {
	return nil, nil
}

type fakeMessages struct {
	rows       []models.Message
	markedUpTo uint
}

func (f *fakeMessages) Create(_ context.Context, m *models.Message) error {
	m.ID = uint(len(f.rows) + 1)
	f.rows = append(f.rows, *m)
	return nil
}

func (f *fakeMessages) ListByMatch(_ context.Context, matchID, beforeID uint, limit int) ([]models.Message, error) {
	var out []models.Message
	for i := len(f.rows) - 1; i >= 0 && len(out) < limit; i-- {
		m := f.rows[i]
		if m.MatchID == matchID && (beforeID == 0 || m.ID < beforeID) {
			out = append(out, m)
		}
	}
	return out, nil
}

func (f *fakeMessages) MarkReadUpTo(_ context.Context, _, _, upToID uint, _ time.Time) error {
	f.markedUpTo = upToID
	return nil
}

type fakeCandidates struct {
	got CandidateQuery
}

func (f *fakeCandidates) FindCandidates(_ context.Context, q CandidateQuery) ([]Candidate, error) {
	f.got = q
	return nil, nil
}

type fakePhotos struct {
	rows []models.UserPhoto
}

func (f *fakePhotos) ListByUser(_ context.Context, userID uint) ([]models.UserPhoto, error) {
	var out []models.UserPhoto
	for _, p := range f.rows {
		if p.UserID == userID {
			out = append(out, p)
		}
	}
	return out, nil
}

func (f *fakePhotos) Create(_ context.Context, p *models.UserPhoto) error {
	p.ID = uint(len(f.rows) + 1)
	f.rows = append(f.rows, *p)
	return nil
}

func (f *fakePhotos) Delete(_ context.Context, userID, photoID uint) (*models.UserPhoto, error) {
	for i, p := range f.rows {
		if p.ID == photoID && p.UserID == userID {
			f.rows = append(f.rows[:i], f.rows[i+1:]...)
			return &p, nil
		}
	}
	return nil, apperr.ErrPhotoNotFound
}

type fakeUpload struct{ storedIn string }

func (f *fakeUpload) Store(dir string) (string, error) {
	f.storedIn = dir
	return dir + "/photo.jpg", nil
}

type fakeFiles struct{ deleted []string }

func (f *fakeFiles) Delete(paths ...string) error {
	f.deleted = append(f.deleted, paths...)
	return nil
}

func user(id uint, gender models.Gender, pref models.Preference, birth time.Time) *models.User {
	u := &models.User{Name: "User", Gender: gender, InterestedIn: pref, BirthDate: birth}
	u.ID = id
	return u
}

func ptr[T any](v T) *T { return &v }
