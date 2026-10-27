package repositories

import (
	"context"
	"time"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"leochat/app/models"
)

type MessageRepository struct {
	orm contractsorm.Orm
}

func NewMessageRepository(o contractsorm.Orm) *MessageRepository {
	return &MessageRepository{orm: o}
}

func (r *MessageRepository) query(ctx context.Context) contractsorm.Query {
	return r.orm.WithContext(ctx).Query()
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.query(ctx).Create(message)
}

func (r *MessageRepository) ListByMatch(ctx context.Context, matchID, beforeID uint, limit int) ([]models.Message, error) {
	q := r.query(ctx).Where("match_id = ?", matchID)
	if beforeID > 0 {
		q = q.Where("id < ?", beforeID)
	}

	var messages []models.Message
	err := q.OrderByDesc("id").Limit(limit).Get(&messages)
	return messages, err
}

func (r *MessageRepository) MarkReadUpTo(ctx context.Context, matchID, readerID, upToID uint, at time.Time) error {
	_, err := r.query(ctx).Model(&models.Message{}).
		Where("match_id = ? AND sender_id <> ? AND id <= ? AND read_at IS NULL", matchID, readerID, upToID).
		Update("read_at", at)
	return err
}
