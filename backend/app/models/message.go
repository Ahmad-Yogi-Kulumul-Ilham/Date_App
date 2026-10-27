package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Message struct {
	orm.Model
	MatchID  uint
	SenderID uint
	Body     string
	ReadAt   *time.Time
}
