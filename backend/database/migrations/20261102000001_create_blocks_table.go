package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261102000001CreateBlocksTable struct{}

func (r *M20261102000001CreateBlocksTable) Signature() string {
	return "20261102000001_create_blocks_table"
}

func (r *M20261102000001CreateBlocksTable) Up() error {
	if facades.Schema().HasTable("blocks") {
		return nil
	}

	return facades.Schema().Create("blocks", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("blocker_id")
		table.UnsignedBigInteger("blocked_id")
		table.Timestamps()

		table.Unique("blocker_id", "blocked_id")
		table.Index("blocked_id")
		table.Foreign("blocker_id").References("id").On("users").CascadeOnDelete()
		table.Foreign("blocked_id").References("id").On("users").CascadeOnDelete()
	})
}

func (r *M20261102000001CreateBlocksTable) Down() error {
	return facades.Schema().DropIfExists("blocks")
}
