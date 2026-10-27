package repositories

import (
	"context"
	"errors"
	"time"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	frameworkerrors "github.com/goravel/framework/errors"

	"leochat/app/apperr"
	"leochat/app/models"
	"leochat/app/services"
)

type MatchRepository struct {
	orm contractsorm.Orm
}

func NewMatchRepository(o contractsorm.Orm) *MatchRepository {
	return &MatchRepository{orm: o}
}

func (r *MatchRepository) query(ctx context.Context) contractsorm.Query {
	return r.orm.WithContext(ctx).Query()
}

func (r *MatchRepository) FirstOrCreate(ctx context.Context, userA, userB uint) (*models.Match, error) {
	one, two := models.OrderedPair(userA, userB)

	var match models.Match
	err := r.query(ctx).
		Where("user_one_id = ? AND user_two_id = ?", one, two).
		FirstOrCreate(&match, models.Match{UserOneID: one, UserTwoID: two})
	if isDuplicateKey(err) {
		// Both users liked each other at the same moment and the other request won.
		match = models.Match{}
		err = r.query(ctx).Where("user_one_id = ? AND user_two_id = ?", one, two).FirstOrFail(&match)
	}
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) FindByID(ctx context.Context, id uint) (*models.Match, error) {
	var match models.Match
	err := r.query(ctx).Where("id = ?", id).FirstOrFail(&match)
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return nil, apperr.ErrMatchNotFound
	}
	if err != nil {
		return nil, err
	}
	return &match, nil
}

const matchSummariesSQL = `
SELECT m.id AS match_id,
       m.created_at AS matched_at,
       u.id AS partner_id,
       u.name AS partner_name,
       u.birth_date AS partner_birth_date,
       (SELECT p.path FROM user_photos p
         WHERE p.user_id = u.id ORDER BY p.position LIMIT 1) AS partner_photo,
       lm.id AS last_message_id,
       lm.sender_id AS last_message_sender_id,
       lm.body AS last_message_body,
       lm.created_at AS last_message_at,
       lm.read_at AS last_message_read_at
FROM matches m
JOIN users u
  ON u.id = IF(m.user_one_id = ?, m.user_two_id, m.user_one_id)
LEFT JOIN messages lm
  ON lm.id = (SELECT MAX(x.id) FROM messages x WHERE x.match_id = m.id)
WHERE m.user_one_id = ? OR m.user_two_id = ?
ORDER BY COALESCE(lm.created_at, m.created_at) DESC
LIMIT 200`

type matchSummaryRow struct {
	MatchID             uint
	MatchedAt           time.Time
	PartnerID           uint
	PartnerName         string
	PartnerBirthDate    time.Time
	PartnerPhoto        *string
	LastMessageID       *uint
	LastMessageSenderID *uint
	LastMessageBody     *string
	LastMessageAt       *time.Time
	LastMessageReadAt   *time.Time
}

func (r *MatchRepository) ListForUser(ctx context.Context, userID uint) ([]services.MatchSummary, error) {
	var rows []matchSummaryRow
	if err := r.query(ctx).Raw(matchSummariesSQL, userID, userID, userID).Scan(&rows); err != nil {
		return nil, err
	}

	summaries := make([]services.MatchSummary, len(rows))
	for i, row := range rows {
		partner := models.User{Name: row.PartnerName, BirthDate: row.PartnerBirthDate}
		partner.ID = row.PartnerID
		if row.PartnerPhoto != nil {
			partner.Photos = []models.UserPhoto{{UserID: row.PartnerID, Path: *row.PartnerPhoto}}
		}

		summary := services.MatchSummary{MatchID: row.MatchID, MatchedAt: row.MatchedAt, Partner: partner}
		if row.LastMessageID != nil {
			summary.LastMessage = &services.MessagePreview{
				ID:       *row.LastMessageID,
				SenderID: *row.LastMessageSenderID,
				Body:     *row.LastMessageBody,
				SentAt:   *row.LastMessageAt,
				ReadAt:   row.LastMessageReadAt,
			}
		}
		summaries[i] = summary
	}
	return summaries, nil
}
