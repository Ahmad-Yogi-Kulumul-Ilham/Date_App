package repositories

import (
	"context"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"leochat/app/apperr"
	"leochat/app/models"
)

type PhotoRepository struct {
	orm contractsorm.Orm
}

func NewPhotoRepository(o contractsorm.Orm) *PhotoRepository {
	return &PhotoRepository{orm: o}
}

func (r *PhotoRepository) query(ctx context.Context) contractsorm.Query {
	return r.orm.WithContext(ctx).Query()
}

func (r *PhotoRepository) ListByUser(ctx context.Context, userID uint) ([]models.UserPhoto, error) {
	var photos []models.UserPhoto
	err := r.query(ctx).Where("user_id = ?", userID).OrderBy("position").Get(&photos)
	return photos, err
}

func (r *PhotoRepository) Create(ctx context.Context, photo *models.UserPhoto) error {
	return r.query(ctx).Create(photo)
}

func (r *PhotoRepository) Delete(ctx context.Context, userID, photoID uint) (*models.UserPhoto, error) {
	var photo models.UserPhoto
	if err := r.query(ctx).Where("id = ? AND user_id = ?", photoID, userID).First(&photo); err != nil {
		return nil, err
	}
	if photo.ID == 0 {
		return nil, apperr.ErrPhotoNotFound
	}

	if _, err := r.query(ctx).Delete(&photo); err != nil {
		return nil, err
	}
	return &photo, nil
}
