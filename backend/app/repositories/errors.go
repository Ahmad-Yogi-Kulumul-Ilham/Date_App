// Package repositories implements the services' persistence ports on top of
// the Goravel ORM and MySQL.
package repositories

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const mysqlDuplicateEntry = 1062

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntry
}
