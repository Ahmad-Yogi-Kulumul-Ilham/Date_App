package repositories

import (
	"context"
	"errors"
	"sort"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/database/orm"
	frameworkerrors "github.com/goravel/framework/errors"

	"leochat/app/apperr"
	"leochat/app/models"
)

type UserRepository struct {
	orm contractsorm.Orm
}

func NewUserRepository(o contractsorm.Orm) *UserRepository {
	return &UserRepository{orm: o}
}

func (r *UserRepository) query(ctx context.Context) contractsorm.Query {
	return r.orm.WithContext(ctx).Query()
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.query(ctx).With("Photos").Where("id = ?", id).FirstOrFail(&user)
	if errors.Is(err, frameworkerrors.OrmRecordNotFound) {
		return nil, apperr.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	sort.SliceStable(user.Photos, func(i, j int) bool {
		return user.Photos[i].Position < user.Photos[j].Position
	})
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.query(ctx).Where("email = ?", email).First(&user); err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	err := r.query(ctx).Create(user)
	if isDuplicateKey(err) {
		return apperr.ErrEmailTaken
	}
	return err
}

func (r *UserRepository) Save(ctx context.Context, user *models.User) error {
	return r.query(ctx).Omit(orm.Associations).Save(user)
}
