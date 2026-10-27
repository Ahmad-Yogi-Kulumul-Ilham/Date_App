package bootstrap

import (
	"github.com/goravel/framework/contracts/database/seeder"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"leochat/app/http/middleware"
	"leochat/config"
	"leochat/database/seeders"
	"leochat/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithSeeders(func() []seeder.Seeder {
			return []seeder.Seeder{&seeders.DemoSeeder{}}
		}).
		WithRouting(func() {
			routes.Web()
			routes.Api()
			routes.Grpc()
		}).
		WithCallback(middleware.RegisterRateLimiters).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
