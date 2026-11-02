package repositories

import (
	"context"

	contractsorm "github.com/goravel/framework/contracts/database/orm"

	"leochat/app/models"
)

type ReportRepository struct {
	orm contractsorm.Orm
}

func NewReportRepository(o contractsorm.Orm) *ReportRepository {
	return &ReportRepository{orm: o}
}

// Create stores the report. A second report on the same user by the same
// reporter is ignored; the first one is already waiting for a moderator.
func (r *ReportRepository) Create(ctx context.Context, report *models.Report) error {
	err := r.orm.WithContext(ctx).Query().Create(report)
	if isDuplicateKey(err) {
		return nil
	}
	return err
}
