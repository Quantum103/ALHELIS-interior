package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfigFromEnv() Config {
	return Config{
		Host:     getEnv("DB_HOST", "postgres_db"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "alhelis"),
		Password: getEnv("DB_PASSWORD", "coffee_password"),
		DBName:   getEnv("DB_NAME", "user_db"),
	}
}

func NewPostgresPool(cfg Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	var currentDB, currentUser, profileName string

	err = pool.QueryRow(
		ctx,
		`
    SELECT current_database(), current_user
    `,
	).Scan(&currentDB, &currentUser)

	if err != nil {
		return nil, fmt.Errorf("diagnostic query: %w", err)
	}

	err = pool.QueryRow(
		ctx,
		`
    SELECT name
    FROM user_profiles
    WHERE user_id = 1
    `,
	).Scan(&profileName)

	if err != nil {
		return nil, fmt.Errorf("diagnostic profile query: %w", err)
	}

	log.Printf(
		"DB DIAGNOSTIC: database=%s user=%s profile=%s",
		currentDB,
		currentUser,
		profileName,
	)

	log.Printf("Успешно подключились к PostgreSQL (%s)", cfg.DBName)

	return pool, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
