package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261026000001CreateUsersTable struct{}

func (r *M20261026000001CreateUsersTable) Signature() string {
	return "20261026000001_create_users_table"
}

func (r *M20261026000001CreateUsersTable) Up() error {
	if facades.Schema().HasTable("users") {
		return nil
	}

	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.ID()
		table.String("name", 50)
		table.String("email")
		table.String("password")
		table.Date("birth_date")
		table.Enum("gender", []any{"male", "female"})
		table.Enum("interested_in", []any{"male", "female", "everyone"})
		table.String("bio", 500).Default("")
		table.Decimal("latitude").Total(10).Places(7).Nullable()
		table.Decimal("longitude").Total(10).Places(7).Nullable()
		table.DateTime("last_active_at").Nullable()
		table.Timestamps()

		table.Unique("email")
		table.Index("latitude", "longitude")
	})
}

func (r *M20261026000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
