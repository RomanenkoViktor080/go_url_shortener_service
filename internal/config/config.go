package config

import (
	"github.com/RomanenkoViktor080/url_shortener_service/internal/util/env"
)

type config struct {
	Port string
	Dns  string
}

func Mount() config {
	return config{
		Port: env.GetString("PORT", "8080"),
		Dns:  env.GetString("GOOSE_DBSTRING", "postgresql://user:password@localhost:5432/postgres"),
	}
}
