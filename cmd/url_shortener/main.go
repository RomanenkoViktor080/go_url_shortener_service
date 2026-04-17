package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/config"
	"github.com/RomanenkoViktor080/url_shortener_service/internal/pkg/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	cfx := config.Mount()

	validator.InitValidatorTranslator()

	pool, err := pgxpool.New(ctx, cfx.Dns)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	slog.Info("connected to the database")

	redis := config.NewRedisClient()
	if err = redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}
	defer redis.Close()

	app := config.Init(cfx, pool, redis)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
