package main

import (
	"context"
	"fmt"
	"hexleturlshort/internal/config"
	"hexleturlshort/internal/database"
	"os"
	"path/filepath"
)

func main() {
	cfg, err := config.ReadFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config from env: %v", err)
		os.Exit(1)
	}

	db, err := database.OpenPostgres(context.Background(), cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open postgres db: %v", err)
		os.Exit(1)
	}

	seedFilePath, err := filepath.Abs(
		filepath.Join("db", "seeds", "development.sql"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build path to seed sql file: %v", err)
		os.Exit(1)
	}

	seedFileContent, err := os.ReadFile(seedFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read seed sql file: %v", err)
		os.Exit(1)
	}

	query := string(seedFileContent)

	result, err := db.ExecContext(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to seed db: %v", err)
		os.Exit(1)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read rows affected: %v", err)
		os.Exit(1)
	}

	fmt.Printf("db seeded: %d rows affected\n", rowsAffected)
}
