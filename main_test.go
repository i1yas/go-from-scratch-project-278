package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/apiapp"
	"hexleturlshort/internal/database/testutils"
)

var (
	binaryPath = ""
	dbURL      = ""
)

func TestSmoke(t *testing.T) {
	port := pickFreePort(t)

	addr := fmt.Sprintf("127.0.0.1:%s", port)
	host := "http://" + addr
	baseURL := "http://example.com"

	t.Run("main routes", func(t *testing.T) {
		ctx := t.Context()

		cmd := exec.CommandContext(ctx, binaryPath)
		cmd.Env = append(cmd.Env,
			"ENV=production",
			"DATABASE_URL="+dbURL,
			"BASE_URL="+baseURL,
			"HTTP_ADDR="+addr,
		)

		err := cmd.Start()
		require.NoError(t, err)

		waitReady(t, host+"/ping", time.Second)

		requireStatus(t, http.StatusOK, host+"/api/links")

		requireStatus(t, http.StatusNotFound, host+"/api/links/999")

		requireStatus(t, http.StatusOK, host+"/api/link_visits")

		requireStatus(t, http.StatusNotFound, host+"/r/unknown-link")
	})

	t.Run("error db down", func(t *testing.T) {
		ctx := t.Context()

		cmd := exec.CommandContext(ctx, binaryPath)
		cmd.Env = append(cmd.Env,
			"ENV=production",
			"DATABASE_URL=postgres://test:test@localhost:9999/unknown-db?sslmode=disable",
			"BASE_URL="+baseURL,
			"HTTP_ADDR="+addr,
		)

		stderr, err := runCommand(cmd)

		require.Error(t, err)
		require.Contains(t, stderr, apiapp.ErrFailedToOpenDB.Error())
	})

	t.Run("failed to start server", func(t *testing.T) {
		ctx := t.Context()

		cmd := exec.CommandContext(ctx, binaryPath)
		cmd.Env = append(cmd.Env,
			"ENV=production",
			"DATABASE_URL="+dbURL,
			"BASE_URL="+baseURL,
			"HTTP_ADDR=invalid",
		)

		var errOut bytes.Buffer

		cmd.Stderr = &errOut

		err := cmd.Run()

		require.Error(t, err)
		require.Contains(t, errOut.String(), apiapp.ErrFailedToStartServer.Error())
	})

	t.Run("invalid sentry dsn", func(t *testing.T) {
		ctx := t.Context()

		cmd := exec.CommandContext(ctx, binaryPath)
		cmd.Env = append(cmd.Env,
			"ENV=production",
			"DATABASE_URL="+dbURL,
			"BASE_URL="+baseURL,
			"HTTP_ADDR="+addr,
			"SENTRY_DSN=invalid",
		)

		stderr, err := runCommand(cmd)

		require.Error(t, err)
		require.Contains(t, stderr, apiapp.ErrFailedToInitSentry.Error())
	})
}

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "hexlet-shorturl-smoketest-*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %s", err.Error())
	}

	buildCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	bin, err := buildTestingBinary(buildCtx, tempDir)
	if err != nil {
		log.Fatalf("Failed to build binary: %s", err.Error())
	}

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Fatalf("Failed to cleanup smoke test dir: %s", err.Error())
		}
	}()

	binaryPath = bin

	pg, err := testutils.SetupTestPostgres(context.Background())
	if err != nil {
		log.Fatalf("Failed to setup test db: %v", err)
	}

	if err := pg.Migrate(filepath.Join("db", "migrations")); err != nil {
		log.Fatalf("Failed to migrate test db: %v", err)
	}

	dbURL = pg.URL

	m.Run()

	if err := pg.Container.Terminate(context.Background()); err != nil {
		log.Fatalf("Failed to terminate test db container: %v", err)
	}
}

func buildTestingBinary(ctx context.Context, dir string) (string, error) {
	bin := filepath.Join(dir, "bin")
	cmd := exec.CommandContext(ctx, "go", "build", "-o", bin, ".")

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return bin, nil
}

func waitReady(t *testing.T, url string, d time.Duration) {
	t.Helper()

	deadline := time.Now().Add(d)

	for time.Now().Before(deadline) {
		if r, err := http.Get(url); err == nil && r.StatusCode == 200 {
			require.NoError(t, r.Body.Close())
			return
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatal("application did not become ready")
}

func runCommand(cmd *exec.Cmd) (string, error) {
	var errOut bytes.Buffer

	cmd.Stderr = &errOut

	err := cmd.Run()

	return errOut.String(), err
}

func requireStatus(t *testing.T, status int, url string) {
	r, err := http.Get(url)

	require.NoError(t, err)
	require.Equal(t, status, r.StatusCode)
}

func pickFreePort(t *testing.T) string {
	t.Helper()

	from := 8080
	to := 8090

	for p := from; p <= to; p++ {
		port := fmt.Sprint(p)

		listener, err := net.Listen("tcp", ":"+port)
		if err == nil {
			require.NoError(t, listener.Close())
			return port
		}
	}

	t.Fatalf("failed to pick free port in range %d-%d", from, to)

	return ""
}
