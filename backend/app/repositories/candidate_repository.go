package repositories

import (
	"context"
	"time"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"leochat/app/models"
	"leochat/app/services"
)

type CandidateRepository struct {
	orm contractsorm.Orm
}

func NewCandidateRepository(o contractsorm.Orm) *CandidateRepository {
	return &CandidateRepository{orm: o}
}

// The bounding-box predicates use the (latitude, longitude) index to narrow the
// scan; ST_Distance_Sphere then gives the exact great-circle distance in meters.
const candidatesSQL = `
SELECT u.id, u.name, u.birth_date, u.gender, u.bio,
       ST_Distance_Sphere(POINT(u.longitude, u.latitude), POINT(?, ?)) AS distance_meters
FROM users u
WHERE u.id <> ?
  AND u.latitude BETWEEN ? AND ?
  AND u.longitude BETWEEN ? AND ?
  AND u.birth_date > ? AND u.birth_date <= ?
  AND u.gender IN ?
  AND u.interested_in IN (?, 'everyone')
  AND NOT EXISTS (
      SELECT 1 FROM swipes s WHERE s.swiper_id = ? AND s.swiped_id = u.id
  )
HAVING distance_meters <= ?
ORDER BY distance_meters
LIMIT ?`

type candidateRow struct {
	ID             uint
	Name           string
	BirthDate      time.Time
	Gender         string
	Bio            string
	DistanceMeters float64
}

func (r *CandidateRepository) FindCandidates(ctx context.Context, q services.CandidateQuery) ([]services.Candidate, error) {
	genders := make([]string, len(q.Genders))
	for i, g := range q.Genders {
		genders[i] = string(g)
	}

	var rows []candidateRow
	err := r.orm.WithContext(ctx).Query().Raw(candidatesSQL,
		q.Origin.Lng, q.Origin.Lat,
		q.ViewerID,
		q.Box.MinLat, q.Box.MaxLat,
		q.Box.MinLng, q.Box.MaxLng,
		q.BornAfter, q.BornOnOrBefore,
		genders,
		string(q.ViewerGender),
		q.ViewerID,
		q.RadiusMeters,
		q.Limit,
	).Scan(&rows)
	if err != nil || len(rows) == 0 {
		return nil, err
	}

	photos, err := r.photosByUser(ctx, rows)
	if err != nil {
		return nil, err
	}

	candidates := make([]services.Candidate, len(rows))
	for i, row := range rows {
		user := models.User{
			Name:      row.Name,
			BirthDate: row.BirthDate,
			Gender:    models.Gender(row.Gender),
			Bio:       row.Bio,
			Photos:    photos[row.ID],
		}
		user.ID = row.ID
		candidates[i] = services.Candidate{User: user, DistanceMeters: row.DistanceMeters}
	}
	return candidates, nil
}

func (r *CandidateRepository) photosByUser(ctx context.Context, rows []candidateRow) (map[uint][]models.UserPhoto, error) {
	ids := make([]any, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}

	var photos []models.UserPhoto
	err := r.orm.WithContext(ctx).Query().
		WhereIn("user_id", ids).
		OrderBy("user_id").
		OrderBy("position").
		Get(&photos)
	if err != nil {
		return nil, err
	}

	byUser := make(map[uint][]models.UserPhoto, len(rows))
	for _, p := range photos {
		byUser[p.UserID] = append(byUser[p.UserID], p)
	}
	return byUser, nil
}
