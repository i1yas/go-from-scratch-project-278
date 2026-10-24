package apiapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/config"
	"hexleturlshort/internal/database/testutils"
)

func TestApp_HappyPath(t *testing.T) {
	dbURL := setupTestDB(t)

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

	gin.SetMode(gin.TestMode)

	step(t, "ping returns pong", func(t *testing.T) {
		w := makeRequest(t, app, "GET", "/ping", nil)
		require.Equal(t, 200, w.Code)

		require.Equal(t, "pong", w.Body.String())
	})

	var createdLink linkDTO

	step(t, "create and get link", func(t *testing.T) {
		url := "http://test.com"
		code := "test"

		w := makeRequest(t, app, "POST", "/api/links", linkDTO{
			OriginalURL: url,
			Shortname:   code,
		})

		require.Equal(t, 201, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)
		require.Equal(t, url, resp.OriginalURL)
		require.Equal(t, code, resp.Shortname)

		createdLink = resp

		w = makeRequest(t, app, "POST", "/api/links", linkDTO{
			OriginalURL: url,
			Shortname:   code,
		})

		require.Equal(t, 409, w.Code)

		w = makeRequest(t, app, "GET", linkURL(resp.ID), nil)

		require.Equal(t, 200, w.Code)

		getByIDResp := bodyAsJSON[linkDTO](t, w)
		require.Equal(t, url, getByIDResp.OriginalURL)
		require.Equal(t, code, getByIDResp.Shortname)
	})

	step(t, "get links", func(t *testing.T) {
		w := makeRequest(t, app, "GET", "/api/links", nil)

		require.Equal(t, 200, w.Code)
		require.Equal(t, "links 0-0/1", w.Header().Get("Content-Range"))

		getLinksResp := bodyAsJSON[[]linkDTO](t, w)
		require.Equal(t, 1, len(getLinksResp))
		require.Equal(t, createdLink, getLinksResp[0])
	})

	step(t, "update link", func(t *testing.T) {
		url := "http://test.com/added-page"
		code := "test-1"

		w := makeRequest(t, app, "PUT", linkURL(createdLink.ID), linkDTO{
			OriginalURL: url,
			Shortname:   code,
		})

		require.Equal(t, 200, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)
		require.Equal(t, url, resp.OriginalURL)
		require.Equal(t, code, resp.Shortname)

		createdLink = resp

		w = makeRequest(t, app, "GET", linkURL(resp.ID), nil)

		require.Equal(t, 200, w.Code)

		getByIDResp := bodyAsJSON[linkDTO](t, w)
		require.Equal(t, url, getByIDResp.OriginalURL)
		require.Equal(t, code, getByIDResp.Shortname)
	})

	step(t, "resolve link", func(t *testing.T) {
		clientIP := "1.2.3.4"
		clientReferer := "http://test-site.com"
		clientUserAgent := "test-user-agent"

		resolveLinkURL := fmt.Sprintf("/r/%s", createdLink.Shortname)
		req, err := http.NewRequest("GET", resolveLinkURL, nil)
		require.NoError(t, err)

		req.Header.Set(gin.PlatformCloudflare, clientIP)
		req.Header.Set("Referer", clientReferer)
		req.Header.Set("User-Agent", clientUserAgent)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		require.Equal(t, 302, w.Code)
		require.Equal(t, createdLink.OriginalURL, w.Header().Get("Location"))

		w = makeRequest(t, app, "GET", "/api/link_visits?range[0,10]", nil)
		require.Equal(t, 200, w.Code)
		require.Equal(t, "visits 0-0/1", w.Header().Get("Content-Range"))

		visitsResp := bodyAsJSON[[]visitDTO](t, w)
		require.Equal(t, 1, len(visitsResp))

		gotVisit := visitsResp[0]

		require.Equal(t, createdLink.ID, gotVisit.LinkID)
		require.Equal(t, clientIP, gotVisit.IP)
		require.Equal(t, clientReferer, gotVisit.Referer)
		require.Equal(t, clientUserAgent, gotVisit.UserAgent)
	})

	step(t, "resolve link without IP, no visit recorded", func(t *testing.T) {
		resolveLinkURL := fmt.Sprintf("/r/%s", createdLink.Shortname)
		req, err := http.NewRequest("GET", resolveLinkURL, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		require.Equal(t, 302, w.Code)
		require.Equal(t, createdLink.OriginalURL, w.Header().Get("Location"))

		w = makeRequest(t, app, "GET", "/api/link_visits?range[0,10]", nil)
		require.Equal(t, 200, w.Code)

		visitsResp := bodyAsJSON[[]visitDTO](t, w)
		require.Equal(t, 1, len(visitsResp))
	})

	step(t, "delete link", func(t *testing.T) {
		w := makeRequest(t, app, "DELETE", linkURL(createdLink.ID), nil)
		require.Equal(t, 204, w.Code)

		w = makeRequest(t, app, "GET", linkURL(createdLink.ID), nil)
		require.Equal(t, 404, w.Code)

		w = makeRequest(t, app, "GET", "/api/link_visits", nil)
		require.Equal(t, 200, w.Code)

		visitsResp := bodyAsJSON[[]visitDTO](t, w)
		require.Equal(t, 0, len(visitsResp))
	})
}

func TestApp_Errors(t *testing.T) {
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
		dbURL := setupTestDB(t)

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

type linkDTO struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	Shortname   string `json:"short_name"`
}

type visitDTO struct {
	ID        int64  `json:"id"`
	LinkID    int64  `json:"link_id"`
	IP        string `json:"ip"`
	Referer   string `json:"reffer"`
	UserAgent string `json:"user_agent"`
}

func linkURL(id int64) string {
	return fmt.Sprintf("/api/links/%d", id)
}

func setupTestDB(t *testing.T) string {
	t.Helper()

	ctx := t.Context()

	pg, err := testutils.SetupTestPostgres(ctx)
	require.NoError(t, err)

	require.NoError(t, pg.Migrate(migrationsDir(t)))

	t.Cleanup(func() {
		err := pg.DB.Close()
		require.NoError(t, err)

		err = pg.Container.Terminate(context.Background())
		require.NoError(t, err)
	})

	return pg.URL
}

func makeRequest(t *testing.T, app *App, method, url string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader io.Reader

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)

		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, req)

	return w
}

func bodyAsJSON[T any](t *testing.T, w *httptest.ResponseRecorder) T {
	t.Helper()

	respRaw := w.Body.Bytes()

	var resp T

	err := json.Unmarshal(respRaw, &resp)
	require.NoError(t, err)

	return resp
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

// step runs a single step in a sequential test scenario
func step(t *testing.T, name string, f func(*testing.T)) {
	t.Helper()

	if !t.Run(name, f) {
		t.FailNow()
	}
}

func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
