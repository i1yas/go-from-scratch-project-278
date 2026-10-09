package apiapp

import (
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"hexleturlshort/internal/config"
)

func TestApp(t *testing.T) {
	dbURL, db := setupTestDB(t)

	seedDB(t, db, "links")
	seedDB(t, db, "visits")

	t.Run("main api routes work", func(t *testing.T) {
		cases := []struct {
			name       string
			method     string
			url        string
			wantStatus int
			wantBody   string
		}{
			{
				name:       "first two links",
				method:     "GET",
				url:        `/api/links?range=[0,1]&sort["id","ASC"]`,
				wantStatus: 200,
				wantBody: `[
					{
						"id": 1,
						"short_name": "google",
						"short_url": "http://example.com/r/google",
						"original_url": "https://google.com"
					},
					{
						"id": 2,
						"short_name": "yandex",
						"short_url": "http://example.com/r/yandex",
						"original_url": "https://yandex.ru"
					}
				]`,
			},
			{
				name:       "link by id",
				method:     "GET",
				url:        `/api/links/1`,
				wantStatus: 200,
				wantBody: `{
					"id": 1,
					"short_name": "google",
					"short_url": "http://example.com/r/google",
					"original_url": "https://google.com"
				}`,
			},
			{
				name:       "link by id not found",
				method:     "GET",
				url:        `/api/links/999`,
				wantStatus: 404,
			},
			{
				name:       "first two visits",
				method:     "GET",
				url:        `/api/link_visits?range=[0,1]&sort["id","ASC"]`,
				wantStatus: 200,
				wantBody: `[
					{
						"id": 1,
						"link_id": 1,
						"created_at": "2026-09-19T07:56:07Z",
						"ip": "36.8.244.30",
						"reffer": "https://example.org/nested/page",
						"user_agent": "very-very-long-user-agent",
						"status": 302
					},
					{
						"id": 2,
						"link_id": 4,
						"created_at": "2026-08-18T10:36:27Z",
						"ip": "144.190.53.110",
						"reffer": "https://example.net/",
						"user_agent": "very-very-long-user-agent",
						"status": 302
					}
				]`,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				ctx := t.Context()

				app, err := New(ctx, config.Config{
					Env: config.EnvProduction,
					App: config.App{
						BaseURL: "http://example.com",
					},
					Database: config.Database{
						URL: dbURL,
					},
					HTTP: config.HTTP{
						Addr: ":8080",
					},
				}, discardLogger())
				require.NoError(t, err)

				req, err := http.NewRequest(tc.method, tc.url, nil)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				app.Router.ServeHTTP(w, req)

				require.Equal(t, tc.wantStatus, w.Code)

				if tc.wantBody != "" {
					require.JSONEq(t, tc.wantBody, w.Body.String())
				}
			})
		}

		t.Run("ping", func(t *testing.T) {
			ctx := t.Context()

			app, err := New(ctx, config.Config{
				Env: config.EnvProduction,
				App: config.App{
					BaseURL: "http://example.com",
				},
				Database: config.Database{
					URL: dbURL,
				},
				HTTP: config.HTTP{
					Addr: ":8080",
				},
			}, discardLogger())
			require.NoError(t, err)

			req, err := http.NewRequest("GET", "/ping", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			require.Equal(t, 200, w.Code)
			require.Equal(t, "pong", w.Body.String())
		})

		t.Run("resolve short link", func(t *testing.T) {
			ctx := t.Context()

			app, err := New(ctx, config.Config{
				Env: config.EnvProduction,
				App: config.App{
					BaseURL: "http://example.com",
				},
				Database: config.Database{
					URL: dbURL,
				},
				HTTP: config.HTTP{
					Addr: ":8080",
				},
			}, discardLogger())
			require.NoError(t, err)

			req, err := http.NewRequest("GET", "/r/google", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			require.Equal(t, 302, w.Code)
			require.Equal(t, "https://google.com", w.Header().Get("Location"))
		})
	})

	t.Run("error no db url", func(t *testing.T) {
		_, err := New(t.Context(), config.Config{
			Env: config.EnvProduction,
			App: config.App{
				BaseURL: "http://example.com",
			},
			HTTP: config.HTTP{
				Addr: ":8080",
			},
		}, discardLogger())
		require.ErrorIs(t, err, ErrFailedToOpenDB)
	})

	t.Run("error db down", func(t *testing.T) {
		_, err := New(t.Context(), config.Config{
			Env: config.EnvProduction,
			Database: config.Database{
				URL: "postgres://test:test@localhost:1234/dbdown?sslmode=disable",
			},
			App: config.App{
				BaseURL: "http://example.com",
			},
			HTTP: config.HTTP{
				Addr: ":8080",
			},
		}, discardLogger())
		require.ErrorIs(t, err, ErrFailedToOpenDB)
	})

	t.Run("error invalid Sentry DSN", func(t *testing.T) {
		_, err := New(t.Context(), config.Config{
			Env: config.EnvProduction,
			Database: config.Database{
				URL: dbURL,
			},
			App: config.App{
				BaseURL: "http://example.com",
			},
			HTTP: config.HTTP{
				Addr: ":8080",
			},
			Sentry: config.Sentry{
				DSN: "not valid",
			},
		}, discardLogger())
		require.ErrorIs(t, err, ErrFailedToInitSentry)
	})
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	m.Run()
}

func setupTestDB(t *testing.T) (string, *sql.DB) {
	t.Helper()

	ctx := t.Context()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("hexlet-shorturl-test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", connString)
	require.NoError(t, err)

	require.NoError(t, db.PingContext(ctx))

	require.NoError(t, goose.SetDialect("postgres"))

	migrations := migrationsDir(t)
	require.NoError(t, goose.Up(db, migrations))

	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	return connString, db
}

func migrationsDir(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to build path to migrations")
	}

	return filepath.Join(
		filename,
		"..", "..", "..",
		"db", "migrations",
	)
}

func seedDB(t *testing.T, db *sql.DB, name string) {
	path := filepath.Join("testdata", "seeds", name+".sql")

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	query := string(data)

	result, err := db.ExecContext(t.Context(), query)
	require.NoError(t, err)

	rows, err := result.RowsAffected()
	require.NoError(t, err)

	require.NotZero(t, rows)
}

func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
