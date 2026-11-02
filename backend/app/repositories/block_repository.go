package repositories

import (
	"context"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"leochat/app/models"
)

type BlockRepository struct {
	orm contractsorm.Orm
}

func NewBlockRepository(o contractsorm.Orm) *BlockRepository {
	return &BlockRepository{orm: o}
}

// Create records the block. Blocking the same user twice is not an error.
func (r *BlockRepository) Create(ctx context.Context, block *models.Block) error {
	err := r.orm.WithContext(ctx).Query().Create(block)
	if isDuplicateKey(err) {
		return nil
	}
	return err
}

func (r *BlockRepository) ExistsBetween(ctx context.Context, userA, userB uint) (bool, error) {
	return r.orm.WithContext(ctx).Query().Model(&models.Block{}).
		Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)", userA, userB, userB, userA).
		Exists()
}
