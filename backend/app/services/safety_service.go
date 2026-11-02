package services

import (
	"context"

	"leochat/app/apperr"
	"leochat/app/models"
)

// SafetyService lets users cut off people they do not want to deal with.
type SafetyService struct {
	users   UserRepository
	blocks  BlockRepository
	matches MatchRepository
}

func NewSafetyService(users UserRepository, blocks BlockRepository, matches MatchRepository) *SafetyService {
	return &SafetyService{users: users, blocks: blocks, matches: matches}
}

// Block hides the two users from each other in discover and ends their match,
// deleting the conversation. Blocking the same user again changes nothing.
func (s *SafetyService) Block(ctx context.Context, blockerID, blockedID uint) error {
	if blockerID == blockedID {
		return apperr.ErrCannotBlockSelf
	}
	if _, err := s.users.FindByID(ctx, blockedID); err != nil {
		return err
	}

	if err := s.blocks.Create(ctx, &models.Block{BlockerID: blockerID, BlockedID: blockedID}); err != nil {
		return err
	}
	return s.matches.DeleteBetween(ctx, blockerID, blockedID)
}
