package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261026000005CreateMessagesTable struct{}

func (r *M20261026000005CreateMessagesTable) Signature() string {
	return "20261026000005_create_messages_table"
}

func (r *M20261026000005CreateMessagesTable) Up() error {
	if facades.Schema().HasTable("messages") {
		return nil
	}

	return facades.Schema().Create("messages", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("match_id")
		table.UnsignedBigInteger("sender_id")
		table.Text("body")
		table.DateTime("read_at").Nullable()
		table.Timestamps()

		table.Index("match_id", "id")
		table.Foreign("match_id").References("id").On("matches").CascadeOnDelete()
		table.Foreign("sender_id").References("id").On("users").CascadeOnDelete()
	})
}

func (r *M20261026000005CreateMessagesTable) Down() error {
	return facades.Schema().DropIfExists("messages")
}
