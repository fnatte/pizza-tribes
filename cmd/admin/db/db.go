package db

import (
	"database/sql"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)


func NewRedisClient() redis.RedisClient {
	return redis.NewRedisClient(&redis.Options{
		Addr:     internal.EnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: internal.EnvOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})
}

func NewSqlClient() *sql.DB {
	sqliteDSN := internal.EnvOrDefault("SQLITE_DSN", "./pizzatribes.db")
	db, err := sql.Open("sqlite3", sqliteDSN)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return db
}
