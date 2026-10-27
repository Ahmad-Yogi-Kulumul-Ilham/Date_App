// Package resources turns domain models into the JSON shapes the API exposes.
// Models are never serialized directly, so internal fields (password hashes,
// exact coordinates) cannot leak by accident.
package resources

import (
	"time"

	"github.com/goravel/framework/support/carbon"

	"leochat/app/geo"
	"leochat/app/models"
	"leochat/app/services"
)

const dateLayout = "2006-01-02"

type Presenter struct {
	photoURL func(path string) string
	now      services.Clock
}

func NewPresenter(photoURL func(path string) string, now services.Clock) *Presenter {
	return &Presenter{photoURL: photoURL, now: now}
}

type Photo struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Position int    `json:"position"`
}

type Me struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	BirthDate    string  `json:"birth_date"`
	Age          int     `json:"age"`
	Gender       string  `json:"gender"`
	InterestedIn string  `json:"interested_in"`
	Bio          string  `json:"bio"`
	HasLocation  bool    `json:"has_location"`
	Photos       []Photo `json:"photos"`
}

type Candidate struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Age        int     `json:"age"`
	Gender     string  `json:"gender"`
	Bio        string  `json:"bio"`
	DistanceKm int     `json:"distance_km"`
	Photos     []Photo `json:"photos"`
}

type Partner struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	PhotoURL *string `json:"photo_url"`
}

type MessagePreview struct {
	Body   string    `json:"body"`
	IsMine bool      `json:"is_mine"`
	SentAt time.Time `json:"sent_at"`
	IsRead bool      `json:"is_read"`
}

type MatchSummary struct {
	ID          uint            `json:"id"`
	MatchedAt   time.Time       `json:"matched_at"`
	Partner     Partner         `json:"partner"`
	LastMessage *MessagePreview `json:"last_message"`
}

type Message struct {
	ID       uint       `json:"id"`
	MatchID  uint       `json:"match_id"`
	SenderID uint       `json:"sender_id"`
	Body     string     `json:"body"`
	IsMine   bool       `json:"is_mine"`
	SentAt   time.Time  `json:"sent_at"`
	ReadAt   *time.Time `json:"read_at"`
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	User        *Me    `json:"user,omitempty"`
}

func (p *Presenter) Me(u *models.User) Me {
	return Me{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		BirthDate:    u.BirthDate.Format(dateLayout),
		Age:          u.AgeAt(p.now()),
		Gender:       string(u.Gender),
		InterestedIn: string(u.InterestedIn),
		Bio:          u.Bio,
		HasLocation:  u.HasLocation(),
		Photos:       p.Photos(u.Photos),
	}
}

func (p *Presenter) Photo(photo models.UserPhoto) Photo {
	return Photo{ID: photo.ID, URL: p.photoURL(photo.Path), Position: photo.Position}
}

func (p *Presenter) Photos(photos []models.UserPhoto) []Photo {
	out := make([]Photo, len(photos))
	for i, photo := range photos {
		out[i] = p.Photo(photo)
	}
	return out
}

func (p *Presenter) Candidates(candidates []services.Candidate) []Candidate {
	now := p.now()
	out := make([]Candidate, len(candidates))
	for i, c := range candidates {
		out[i] = Candidate{
			ID:         c.User.ID,
			Name:       c.User.Name,
			Age:        c.User.AgeAt(now),
			Gender:     string(c.User.Gender),
			Bio:        c.User.Bio,
			DistanceKm: geo.DisplayKm(c.DistanceMeters),
			Photos:     p.Photos(c.User.Photos),
		}
	}
	return out
}

func (p *Presenter) MatchSummaries(summaries []services.MatchSummary, viewerID uint) []MatchSummary {
	now := p.now()
	out := make([]MatchSummary, len(summaries))
	for i, s := range summaries {
		partner := Partner{ID: s.Partner.ID, Name: s.Partner.Name, Age: s.Partner.AgeAt(now)}
		if len(s.Partner.Photos) > 0 {
			url := p.photoURL(s.Partner.Photos[0].Path)
			partner.PhotoURL = &url
		}

		summary := MatchSummary{ID: s.MatchID, MatchedAt: s.MatchedAt, Partner: partner}
		if m := s.LastMessage; m != nil {
			summary.LastMessage = &MessagePreview{
				Body:   m.Body,
				IsMine: m.SenderID == viewerID,
				SentAt: m.SentAt,
				IsRead: m.ReadAt != nil,
			}
		}
		out[i] = summary
	}
	return out
}

func (p *Presenter) Message(m *models.Message, viewerID uint) Message {
	return Message{
		ID:       m.ID,
		MatchID:  m.MatchID,
		SenderID: m.SenderID,
		Body:     m.Body,
		IsMine:   m.SenderID == viewerID,
		SentAt:   stdTime(m.CreatedAt),
		ReadAt:   m.ReadAt,
	}
}

func (p *Presenter) Messages(messages []models.Message, viewerID uint) []Message {
	out := make([]Message, len(messages))
	for i := range messages {
		out[i] = p.Message(&messages[i], viewerID)
	}
	return out
}

func stdTime(t *carbon.DateTime) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.StdTime()
}
