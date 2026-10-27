package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261026000002CreateUserPhotosTable struct{}

func (r *M20261026000002CreateUserPhotosTable) Signature() string {
	return "20261026000002_create_user_photos_table"
}

func (r *M20261026000002CreateUserPhotosTable) Up() error {
	if facades.Schema().HasTable("user_photos") {
		return nil
	}

	return facades.Schema().Create("user_photos", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_id")
		table.String("path")
		table.UnsignedTinyInteger("position").Default(0)
		table.Timestamps()

		table.Index("user_id", "position")
		table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
	})
}

func (r *M20261026000002CreateUserPhotosTable) Down() error {
	return facades.Schema().DropIfExists("user_photos")
}
