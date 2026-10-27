package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261026000004CreateMatchesTable struct{}

func (r *M20261026000004CreateMatchesTable) Signature() string {
	return "20261026000004_create_matches_table"
}

func (r *M20261026000004CreateMatchesTable) Up() error {
	if facades.Schema().HasTable("matches") {
		return nil
	}

	return facades.Schema().Create("matches", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("user_one_id")
		table.UnsignedBigInteger("user_two_id")
		table.Timestamps()

		table.Unique("user_one_id", "user_two_id")
		table.Index("user_two_id")
		table.Foreign("user_one_id").References("id").On("users").CascadeOnDelete()
		table.Foreign("user_two_id").References("id").On("users").CascadeOnDelete()
	})
}

func (r *M20261026000004CreateMatchesTable) Down() error {
	return facades.Schema().DropIfExists("matches")
}
