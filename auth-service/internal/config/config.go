package config

import (
	"os"
)

type Config struct {
	JWTSecret []byte
}

func Load() *Config {
	return &Config{
		JWTSecret: []byte(getEnv("JWT_SECRET", "default-secret-if-not-set")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
