package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrFailedToReadDotenv = errors.New("failed to read .env file")
	ErrUnsupportedEnv     = errors.New("unsupported environment")
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

// ReadFromEnv reads config from environment and .env file
func ReadFromEnv() (Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("%w: %w",
			ErrFailedToReadDotenv, err,
		)
	}

	env, err := readEnvConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Env:      env,
		App:      readAppConfig(),
		Database: readDBConfig(),
		HTTP:     readHTTPConfig(),
		Sentry:   readSentryConfig(),
	}, nil
}

func readEnvConfig() (Environment, error) {
	env := Environment(os.Getenv("ENV"))

	switch env {
	case Environment(""):
		return EnvDevelopment, nil
	case EnvDevelopment, EnvProduction:
		return env, nil
	default:
		return Environment(""), fmt.Errorf("%w: '%s'",
			ErrUnsupportedEnv, env,
		)
	}
}

func readAppConfig() App {
	return App{
		BaseURL: os.Getenv("BASE_URL"),
	}
}

func readDBConfig() Database {
	return Database{
		URL: os.Getenv("DATABASE_URL"),
	}
}

func readHTTPConfig() HTTP {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return HTTP{
		Addr: addr,
	}
}

func readSentryConfig() Sentry {
	return Sentry{
		DSN: os.Getenv("SENTRY_DSN"),
	}
}
