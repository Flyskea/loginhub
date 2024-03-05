package data

import (
	"errors"
	"loginhub/internal/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	mysqlDriver = "mysql"
)

func NewDB(conf *conf.Database) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	switch conf.Driver {
	case mysqlDriver:
		db, err = gorm.Open(mysql.Open(conf.Source), &gorm.Config{})
	default:
		return nil, errors.New("unsupported driver")
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}
