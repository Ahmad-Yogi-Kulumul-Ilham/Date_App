package models

import "github.com/goravel/framework/database/orm"

// Match links two users who liked each other. The pair is stored ordered
// (UserOneID < UserTwoID) so a unique index can prevent duplicates.
type Match struct {
	orm.Model
	UserOneID uint
	UserTwoID uint
}

// OrderedPair returns a and b sorted ascending, the form a Match is stored in.
func OrderedPair(a, b uint) (uint, uint) {
	if a < b {
		return a, b
	}
	return b, a
}

func (m *Match) Includes(userID uint) bool {
	return m.UserOneID == userID || m.UserTwoID == userID
}

// PartnerOf returns the other participant. The caller must ensure Includes(userID).
func (m *Match) PartnerOf(userID uint) uint {
	if m.UserOneID == userID {
		return m.UserTwoID
	}
	return m.UserOneID
}
