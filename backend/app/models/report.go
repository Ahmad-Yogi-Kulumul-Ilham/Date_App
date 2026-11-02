package models

import "github.com/goravel/framework/database/orm"

type ReportReason string

const (
	ReportFakeProfile   ReportReason = "fake_profile"
	ReportInappropriate ReportReason = "inappropriate"
	ReportHarassment    ReportReason = "harassment"
	ReportSpam          ReportReason = "spam"
	ReportUnderage      ReportReason = "underage"
	ReportOther         ReportReason = "other"
)

// Report flags a user for moderators. Each reporter can report a given user once.
type Report struct {
	orm.Model
	ReporterID uint
	ReportedID uint
	Reason     ReportReason
	Details    string
}
