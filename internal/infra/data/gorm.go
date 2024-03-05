package data

import (
	"github.com/Flyskea/gotools/errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"loginhub/internal/base/reason"
)

const (
	ErrMySQLDupEntry            = 1062
	ErrMySQLDupEntryWithKeyName = 1586
)

func defaultErr(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return errors.NotFound(reason.ObjectNotFoundError)
	default:
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
}

func GromToError(err error) error {
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *mysql.MySQLError:
		switch e.Number {
		case ErrMySQLDupEntry, ErrMySQLDupEntryWithKeyName:
			return errors.Conflict(reason.DatabaseError)
		default:
			return defaultErr(err)
		}
	default:
		return defaultErr(err)
	}
}
