package models

import "github.com/goravel/framework/database/orm"

const MaxPhotosPerUser = 6

type UserPhoto struct {
	orm.Model
	UserID   uint
	Path     string
	Position int
}
