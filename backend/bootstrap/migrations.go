package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20261026000001CreateUsersTable{},
		&migrations.M20261026000002CreateUserPhotosTable{},
		&migrations.M20261026000003CreateSwipesTable{},
		&migrations.M20261026000004CreateMatchesTable{},
		&migrations.M20261026000005CreateMessagesTable{},
	}
}
