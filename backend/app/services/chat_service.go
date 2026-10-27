package services

import (
	"context"
	"strings"

	"leochat/app/apperr"
	"leochat/app/models"
)

const (
	DefaultMessagePageSize = 30
	MaxMessagePageSize     = 100
)

type ChatService struct {
	matches  MatchRepository
	messages MessageRepository
	now      Clock
}

func NewChatService(matches MatchRepository, messages MessageRepository, now Clock) *ChatService {
	return &ChatService{matches: matches, messages: messages, now: now}
}

func (s *ChatService) ListMatches(ctx context.Context, userID uint) ([]MatchSummary, error) {
	return s.matches.ListForUser(ctx, userID)
}

// Messages returns a page of messages, newest first. Loading the newest page
// also marks the partner's messages in it as read.
func (s *ChatService) Messages(ctx context.Context, userID, matchID, beforeID uint, limit int) ([]models.Message, error) {
	if _, err := s.participantMatch(ctx, userID, matchID); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = DefaultMessagePageSize
	}
	limit = min(limit, MaxMessagePageSize)

	messages, err := s.messages.ListByMatch(ctx, matchID, beforeID, limit)
	if err != nil {
		return nil, err
	}

	if beforeID == 0 && len(messages) > 0 {
		if err := s.messages.MarkReadUpTo(ctx, matchID, userID, messages[0].ID, s.now()); err != nil {
			return nil, err
		}
	}
	return messages, nil
}

func (s *ChatService) Send(ctx context.Context, userID, matchID uint, body string) (*models.Message, error) {
	if _, err := s.participantMatch(ctx, userID, matchID); err != nil {
		return nil, err
	}

	body = strings.TrimSpace(body)
	if body == "" {
		return nil, apperr.ErrMessageEmpty
	}

	message := &models.Message{MatchID: matchID, SenderID: userID, Body: body}
	if err := s.messages.Create(ctx, message); err != nil {
		return nil, err
	}
	return message, nil
}

// participantMatch loads the match and hides it from non-participants with a
// not-found error, so match IDs cannot be probed.
func (s *ChatService) participantMatch(ctx context.Context, userID, matchID uint) (*models.Match, error) {
	match, err := s.matches.FindByID(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if !match.Includes(userID) {
		return nil, apperr.ErrMatchNotFound
	}
	return match, nil
}
