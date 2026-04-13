package config

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Application struct {
	config config
}
type config struct {
	port     string
	dbConfig dbConfig
}
type dbConfig struct {
	dns string
}

func Init() Application {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initLogger()

	cfg := config{
		port: getEnv("PORT", "8080"),
	}
	return Application{
		config: cfg,
	}
}

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func (app *Application) Run() error {
	handler := app.mount()
	srv := &http.Server{
		Addr:         ":" + app.config.port,
		Handler:      handler,
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Listening on port %s", app.config.port)

	return srv.ListenAndServe()
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
