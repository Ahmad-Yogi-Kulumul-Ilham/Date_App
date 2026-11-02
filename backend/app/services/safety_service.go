package services

import (
	"context"
	"strings"

	"leochat/app/apperr"
	"leochat/app/models"
)

// SafetyService lets users cut off people they do not want to deal with.
type SafetyService struct {
	users   UserRepository
	blocks  BlockRepository
	reports ReportRepository
	matches MatchRepository
}

func NewSafetyService(users UserRepository, blocks BlockRepository, reports ReportRepository, matches MatchRepository) *SafetyService {
	return &SafetyService{users: users, blocks: blocks, reports: reports, matches: matches}
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
	return s.block(ctx, blockerID, blockedID)
}

// Report stores the report for moderators and also blocks the reported user,
// so the reporter does not need a second step to get away from them.
func (s *SafetyService) Report(ctx context.Context, reporterID, reportedID uint, reason models.ReportReason, details string) error {
	if reporterID == reportedID {
		return apperr.ErrCannotReportSelf
	}
	if _, err := s.users.FindByID(ctx, reportedID); err != nil {
		return err
	}

	report := &models.Report{
		ReporterID: reporterID,
		ReportedID: reportedID,
		Reason:     reason,
		Details:    strings.TrimSpace(details),
	}
	if err := s.reports.Create(ctx, report); err != nil {
		return err
	}
	return s.block(ctx, reporterID, reportedID)
}

func (s *SafetyService) block(ctx context.Context, blockerID, blockedID uint) error {
	if err := s.blocks.Create(ctx, &models.Block{BlockerID: blockerID, BlockedID: blockedID}); err != nil {
		return err
	}
	return s.matches.DeleteBetween(ctx, blockerID, blockedID)
}
