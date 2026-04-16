package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		val, err := strconv.Atoi(val)
		if err != nil {
			return fallback
		}
		return val
	}
	return fallback
}
