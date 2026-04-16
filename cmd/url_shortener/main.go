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
	cxt := context.Background()

	validator.InitValidatorTranslator()

	cfx := config.Mount()
	pool, err := pgxpool.New(cxt, cfx.Dns)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	app := config.Init(cfx, pool)

	slog.Info("connected to the database")

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
