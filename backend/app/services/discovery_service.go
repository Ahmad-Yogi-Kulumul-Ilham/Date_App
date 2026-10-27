package services

import (
	"context"
	"time"

	"leochat/app/apperr"
	"leochat/app/geo"
)

const (
	DefaultRadiusKm = 25
	MaxRadiusKm     = 150
	DefaultMaxAge   = 60
	MaxAgeLimit     = 99
	DefaultPageSize = 20
	MaxPageSize     = 50
)

type DiscoveryService struct {
	users      UserRepository
	candidates CandidateRepository
	now        Clock
}

func NewDiscoveryService(users UserRepository, candidates CandidateRepository, now Clock) *DiscoveryService {
	return &DiscoveryService{users: users, candidates: candidates, now: now}
}

// DiscoverFilter comes from the client; zero values mean "use the default".
type DiscoverFilter struct {
	RadiusKm int
	MinAge   int
	MaxAge   int
	Limit    int
}

func (f DiscoverFilter) Normalized() DiscoverFilter {
	out := f
	if out.RadiusKm <= 0 {
		out.RadiusKm = DefaultRadiusKm
	}
	out.RadiusKm = min(out.RadiusKm, MaxRadiusKm)

	if out.MinAge < MinimumAge {
		out.MinAge = MinimumAge
	}
	if out.MaxAge <= 0 {
		out.MaxAge = DefaultMaxAge
	}
	out.MaxAge = min(out.MaxAge, MaxAgeLimit)
	if out.MaxAge < out.MinAge {
		out.MaxAge = out.MinAge
	}

	if out.Limit <= 0 {
		out.Limit = DefaultPageSize
	}
	out.Limit = min(out.Limit, MaxPageSize)
	return out
}

// Discover returns nearby people the viewer has not swiped yet, closest first.
// Both sides' preferences must match: the viewer wants the candidate's gender and
// the candidate wants the viewer's.
func (s *DiscoveryService) Discover(ctx context.Context, viewerID uint, filter DiscoverFilter) ([]Candidate, error) {
	viewer, err := s.users.FindByID(ctx, viewerID)
	if err != nil {
		return nil, err
	}
	if !viewer.HasLocation() {
		return nil, apperr.ErrLocationRequired
	}

	f := filter.Normalized()
	origin := geo.Point{Lat: *viewer.Latitude, Lng: *viewer.Longitude}
	now := s.now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return s.candidates.FindCandidates(ctx, CandidateQuery{
		ViewerID:     viewerID,
		Origin:       origin,
		Box:          geo.BoundingBoxAround(origin, float64(f.RadiusKm)),
		RadiusMeters: float64(f.RadiusKm) * 1000,
		// Age <= MaxAge  <=>  born after (today - (MaxAge+1) years).
		BornAfter: today.AddDate(-(f.MaxAge + 1), 0, 0),
		// Age >= MinAge  <=>  born on or before (today - MinAge years).
		BornOnOrBefore: today.AddDate(-f.MinAge, 0, 0),
		Genders:        viewer.InterestedIn.Genders(),
		ViewerGender:   viewer.Gender,
		Limit:          f.Limit,
	})
}
