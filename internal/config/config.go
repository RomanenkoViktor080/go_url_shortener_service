package config

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	Config       config
	dbConnection *pgxpool.Pool
}
type config struct {
	Port string
	Dns  string
}

func Mount() config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	return config{
		Port: env.GetString("PORT", "8080"),
		Dns:  env.GetString("GOOSE_DBSTRING", "postgresql://user:password@localhost:5432/postgres"),
	}
}
func Init(cfg config, connection *pgxpool.Pool) application {
	initLogger()

	return application{
		Config:       cfg,
		dbConnection: connection,
	}
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func (app *application) Run() error {
	handler := app.mount()
	srv := &http.Server{
		Addr:         ":" + app.Config.Port,
		Handler:      handler,
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
		IdleTimeout:  time.Minute,
	}

	slog.Info("server started", "port", app.Config.Port)

	return srv.ListenAndServe()
}
