package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"leochat/app/facades"
)

type M20261102000002CreateReportsTable struct{}

func (r *M20261102000002CreateReportsTable) Signature() string {
	return "20261102000002_create_reports_table"
}

func (r *M20261102000002CreateReportsTable) Up() error {
	if facades.Schema().HasTable("reports") {
		return nil
	}

	return facades.Schema().Create("reports", func(table schema.Blueprint) {
		table.ID()
		table.UnsignedBigInteger("reporter_id")
		table.UnsignedBigInteger("reported_id")
		table.Enum("reason", []any{"fake_profile", "inappropriate", "harassment", "spam", "underage", "other"})
		table.String("details", 500).Default("")
		table.Timestamps()

		table.Unique("reporter_id", "reported_id")
		table.Index("reported_id")
		table.Foreign("reporter_id").References("id").On("users").CascadeOnDelete()
		table.Foreign("reported_id").References("id").On("users").CascadeOnDelete()
	})
}

func (r *M20261102000002CreateReportsTable) Down() error {
	return facades.Schema().DropIfExists("reports")
}
