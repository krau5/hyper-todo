package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CookieDomain     string
	JwtSecretKey     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
}

func loadConfig() *Config {
	godotenv.Load()

	return &Config{
		CookieDomain:     getEnv("COOKIE_DOMAIN", "localhost"),
		JwtSecretKey:     getEnv("JWT_SECRET_KEY", "secret"),
		PostgresUser:     getEnv("POSTGRES_USER", "user"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:       getEnv("POSTGRES_DB", "hypertodo"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
	}
}

var Envs = loadConfig()

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func GetDsn() string {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:5432/%s",
		Envs.PostgresUser,
		Envs.PostgresPassword,
		Envs.PostgresHost,
		Envs.PostgresDB,
	)

	return dsn
}
