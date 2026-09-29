package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
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

// Config contains app configuration
type Config struct {
	App      App
	Database Database
	HTTP     HTTP
}

// TODO: validate configs

// ReadFromEnv reads config from environment and .env file
func ReadFromEnv() (Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	cfg := Config{
		App:      readAppConfig(),
		Database: readDBConfig(),
		HTTP:     readHTTPConfig(),
	}

	return cfg, nil
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
