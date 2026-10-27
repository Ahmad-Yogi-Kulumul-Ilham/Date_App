package services

import (
	"context"

	"leochat/app/apperr"
	"leochat/app/models"
)

type SwipeService struct {
	users   UserRepository
	swipes  SwipeRepository
	matches MatchRepository
}

func NewSwipeService(users UserRepository, swipes SwipeRepository, matches MatchRepository) *SwipeService {
	return &SwipeService{users: users, swipes: swipes, matches: matches}
}

// SwipeResult carries the new match when the swipe completed a mutual like.
type SwipeResult struct {
	Match *models.Match
}

func (r SwipeResult) Matched() bool { return r.Match != nil }

func (s *SwipeService) Swipe(ctx context.Context, swiperID, swipedID uint, direction models.SwipeDirection) (SwipeResult, error) {
	if swiperID == swipedID {
		return SwipeResult{}, apperr.ErrCannotSwipeSelf
	}
	if _, err := s.users.FindByID(ctx, swipedID); err != nil {
		return SwipeResult{}, err
	}

	already, err := s.swipes.Exists(ctx, swiperID, swipedID)
	if err != nil {
		return SwipeResult{}, err
	}
	if already {
		return SwipeResult{}, apperr.ErrAlreadySwiped
	}

	swipe := &models.Swipe{SwiperID: swiperID, SwipedID: swipedID, Direction: direction}
	if err := s.swipes.Create(ctx, swipe); err != nil {
		return SwipeResult{}, err
	}

	if direction != models.SwipeLike {
		return SwipeResult{}, nil
	}

	likedBack, err := s.swipes.HasLiked(ctx, swipedID, swiperID)
	if err != nil || !likedBack {
		return SwipeResult{}, err
	}

	match, err := s.matches.FirstOrCreate(ctx, swiperID, swipedID)
	if err != nil {
		return SwipeResult{}, err
	}
	return SwipeResult{Match: match}, nil
}
