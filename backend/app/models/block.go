package models

import "github.com/goravel/framework/database/orm"

// Block hides two users from each other in both directions. It is stored once,
// from the blocker's side; the blocked user is never told.
type Block struct {
	orm.Model
	BlockerID uint
	BlockedID uint
}
