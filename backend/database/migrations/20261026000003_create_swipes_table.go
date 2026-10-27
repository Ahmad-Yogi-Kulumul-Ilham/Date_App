package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261026000003CreateSwipesTable struct{}

func (r *M20261026000003CreateSwipesTable) Signature() string {
	return "20261026000003_create_swipes_table"
}

func (r *M20261026000003CreateSwipesTable) Up() error {
	if facades.Schema().HasTable("swipes") {
		return nil
	}

	return facades.Schema().Create("swipes", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("swiper_id")
		table.UnsignedBigInteger("swiped_id")
		table.Enum("direction", []any{"like", "pass"})
		table.Timestamps()

		table.Unique("swiper_id", "swiped_id")
		table.Index("swiped_id")
		table.Foreign("swiper_id").References("id").On("users").CascadeOnDelete()
		table.Foreign("swiped_id").References("id").On("users").CascadeOnDelete()
	})
}

func (r *M20261026000003CreateSwipesTable) Down() error {
	return facades.Schema().DropIfExists("swipes")
}
