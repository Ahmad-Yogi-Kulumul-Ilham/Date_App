package models

import "github.com/goravel/framework/database/orm"

type SwipeDirection string

const (
	SwipeLike SwipeDirection = "like"
	SwipePass SwipeDirection = "pass"
)

type Swipe struct {
	orm.Model
	SwiperID  uint
	SwipedID  uint
	Direction SwipeDirection
}
