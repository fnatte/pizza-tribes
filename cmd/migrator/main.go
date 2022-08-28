package main

import (
	"context"
	"embed"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func ensureWorld(ctx context.Context, r redis.RedisClient) error {
	world := game.NewWorldService(r)
	if err := world.Initialize(ctx); err != nil {
		return err
	}

	return nil
}

func redisMigrations(ctx context.Context) {
	r := redis.NewRedisClient(&redis.Options{
		Addr:     internal.EnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: internal.EnvOrDefault("REDIS_PASSWORD", ""),
		DB:       0,  // use default DB
	})

	err := ensureWorld(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure world")
	}
}

func sqlMigrations(ctx context.Context) {
	sqliteDSN := internal.EnvOrDefault("SQLITE_DSN", "./pizzatribes.db")
	db, err := goose.OpenDBWithDriver("sqlite3", sqliteDSN)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	goose.SetBaseFS(embedMigrations)

    if err := goose.Up(db, "migrations"); err != nil {
        panic(err)
    }
}

func main() {
	log.Info().Msg("Starting migrator")

	ctx := context.Background()
	redisMigrations(ctx)
	sqlMigrations(ctx)

	log.Info().Msg("Migrator done")
}
