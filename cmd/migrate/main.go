package main

import (
	"context"
	"database/sql"
	"loginhub/internal/conf"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	rootCmd.Flags().StringVarP(&flagconf, "conf", "c", "./configs/config-migrate-example.yaml", "config path")
}

func runMigrate() {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	conf := conf.Migrate{}
	err := c.Load()
	if err != nil {
		panic(err)
	}
	err = c.Scan(&conf)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(conf.GetDatabase().GetDriver(), conf.GetDatabase().GetSource())
	if err != nil {
		panic(err)
	}

	err = goose.SetDialect(conf.GetDatabase().Driver)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = goose.VersionContext(ctx, db, conf.GetDir())
	if err != nil {
		panic(err)
	}
	err = goose.UpContext(ctx, db, conf.GetDir())
	if err != nil {
		panic(err)
	}
}

func main() {
	Execute()
}
