package data

import (
	"context"

	"github.com/redis/go-redis/v9"

	"loginhub/internal/conf"
)

func NewRDB(conf *conf.Redis) (redis.UniversalClient, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		WriteTimeout: conf.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.ReadTimeout.AsDuration(),
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := rdb.Close(); err != nil {
			panic(err)
		}
	}
	return rdb, cleanup, nil
}
