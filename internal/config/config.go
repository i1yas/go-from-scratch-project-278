package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Environment of application
type Environment string

const (
	// EnvProduction is production environment
	EnvProduction = Environment("production")
	// EnvDevelopment is development environment
	EnvDevelopment = Environment("development")
)

// App is application specific config
type App struct {
	BaseURL string
}

// Database is database specific config
type Database struct {
	URL string
}

// HTTP is http specific config
type HTTP struct {
	Addr string
}

// Sentry is sentry specific config
type Sentry struct {
	DSN string
}

// Config contains app configuration
type Config struct {
	Env      Environment
	App      App
	Database Database
	HTTP     HTTP
	Sentry   Sentry
}

// TODO: validate configs

// ReadFromEnv reads config from environment and .env file
func ReadFromEnv() (Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	cfg := Config{
		Env:      readEnvConfig(EnvDevelopment),
		App:      readAppConfig(),
		Database: readDBConfig(),
		HTTP:     readHTTPConfig(),
		Sentry:   readSentryConfig(),
	}

	return cfg, nil
}

func readEnvConfig(fallback Environment) Environment {
	env := Environment(os.Getenv("ENV"))

	switch env {
	case EnvDevelopment, EnvProduction:
		return env
	default:
		return fallback
	}
}

func readAppConfig() App {
	return App{
		BaseURL: os.Getenv("APP_BASE_URL"),
	}
}

func readDBConfig() Database {
	return Database{
		URL: os.Getenv("DATABASE_URL"),
	}
}

func readHTTPConfig() HTTP {
	return HTTP{
		Addr: os.Getenv("HTTP_ADDR"),
	}
}

func readSentryConfig() Sentry {
	return Sentry{
		DSN: os.Getenv("SENTRY_DSN"),
	}
}
