package repositories

import (
	"context"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"leochat/app/apperr"
	"leochat/app/models"
)

type SwipeRepository struct {
	orm contractsorm.Orm
}

func NewSwipeRepository(o contractsorm.Orm) *SwipeRepository {
	return &SwipeRepository{orm: o}
}

func (r *SwipeRepository) query(ctx context.Context) contractsorm.Query {
	return r.orm.WithContext(ctx).Query().Model(&models.Swipe{})
}

func (r *SwipeRepository) Exists(ctx context.Context, swiperID, swipedID uint) (bool, error) {
	return r.query(ctx).
		Where("swiper_id = ? AND swiped_id = ?", swiperID, swipedID).
		Exists()
}

func (r *SwipeRepository) HasLiked(ctx context.Context, swiperID, swipedID uint) (bool, error) {
	return r.query(ctx).
		Where("swiper_id = ? AND swiped_id = ? AND direction = ?", swiperID, swipedID, models.SwipeLike).
		Exists()
}

func (r *SwipeRepository) Create(ctx context.Context, swipe *models.Swipe) error {
	err := r.orm.WithContext(ctx).Query().Create(swipe)
	if isDuplicateKey(err) {
		return apperr.ErrAlreadySwiped
	}
	return err
}
