package apiapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

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
		//nolint:lll
		clientReferer := "http://example.com/page/page/page?param1=29a6e0c3-d99d-47e7-b85d-57182cc2d4fe&param2=5f591642-c757-45a1-99a6-31bda9d02a77"
		//nolint:lll
		clientUserAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.15"

		resolveLinkURL := fmt.Sprintf("/r/%s", createdLink.Shortname)
		req, err := http.NewRequest("GET", resolveLinkURL, nil)
		require.NoError(t, err)

		req.Header.Set(gin.PlatformCloudflare, clientIP)
		req.Header.Set("Referer", clientReferer)
		req.Header.Set("User-Agent", clientUserAgent)

		step(t, "first visit", func(t *testing.T) {
			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			require.Equal(t, 302, w.Code)
			require.Equal(t, createdLink.OriginalURL, w.Header().Get("Location"))

			w = makeRequest(t, app, "GET", "/api/link_visits?range=[0,0]", nil)
			require.Equal(t, 200, w.Code)
			require.Equal(t, "link_visits 0-0/1", w.Header().Get("Content-Range"))

			visitsResp := bodyAsJSON[[]visitDTO](t, w)
			require.Equal(t, 1, len(visitsResp))

			gotVisit := visitsResp[0]

			require.Equal(t, createdLink.ID, gotVisit.LinkID)
			require.Equal(t, clientIP, gotVisit.IP)
			require.Equal(t, clientReferer, gotVisit.Referer)
			require.Equal(t, clientUserAgent, gotVisit.UserAgent)
		})

		step(t, "second visit", func(t *testing.T) {
			w := httptest.NewRecorder()
			app.Router.ServeHTTP(w, req)

			require.Equal(t, 302, w.Code)
			require.Equal(t, createdLink.OriginalURL, w.Header().Get("Location"))

			w = makeRequest(t, app, "GET", "/api/link_visits?range=[1,1]", nil)
			require.Equal(t, 200, w.Code)
			require.Equal(t, "link_visits 1-1/2", w.Header().Get("Content-Range"))

			visitsResp := bodyAsJSON[[]visitDTO](t, w)
			require.Equal(t, 1, len(visitsResp))

			gotVisit := visitsResp[0]

			require.Equal(t, createdLink.ID, gotVisit.LinkID)
			require.Equal(t, clientIP, gotVisit.IP)
			require.Equal(t, clientReferer, gotVisit.Referer)
			require.Equal(t, clientUserAgent, gotVisit.UserAgent)
		})
	})

	step(t, "resolve link without IP, no visit recorded", func(t *testing.T) {
		totalVisitsBefore := fetchVisitsTotalCount(t, app)

		resolveLinkURL := fmt.Sprintf("/r/%s", createdLink.Shortname)
		req, err := http.NewRequest("GET", resolveLinkURL, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		require.Equal(t, 302, w.Code)
		require.Equal(t, createdLink.OriginalURL, w.Header().Get("Location"))

		totalVisitsAfter := fetchVisitsTotalCount(t, app)
		totalVisitsInc := totalVisitsAfter - totalVisitsBefore
		require.Equal(t, 0, totalVisitsInc)
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

func TestApp_RequestBigBody(t *testing.T) {
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

	var createdLink linkDTO

	step(t, "create link", func(t *testing.T) {
		urlBase := "http://x?"
		symbolRepeats := 200 - len(urlBase)
		// NOTE: json encode symbol &
		url := urlBase + strings.Repeat("\u0026", symbolRepeats)

		// NOTE: 𝒜 - math script symbol, takes 4 bytes
		code := strings.Repeat("𝒜", 32)

		w := makeRequest(t, app, "POST", "/api/links", linkDTO{
			OriginalURL: url,
			Shortname:   code,
		})

		require.Equal(t, 201, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)
		require.Equal(t, url, resp.OriginalURL)
		require.Equal(t, code, resp.Shortname)

		createdLink = resp
	})

	step(t, "update link", func(t *testing.T) {
		urlBase := "http://y?"
		symbolRepeats := 200 - len(urlBase)
		url := urlBase + strings.Repeat("\u0026", symbolRepeats)

		code := strings.Repeat("𝔅", 32)

		require.GreaterOrEqual(t, len(code), 32)

		w := makeRequest(t, app, "PUT", linkURL(createdLink.ID), linkDTO{
			OriginalURL: url,
			Shortname:   code,
		})

		require.Equal(t, 200, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)
		require.Equal(t, url, resp.OriginalURL)
		require.Equal(t, code, resp.Shortname)
	})

	step(t, "body is too big", func(t *testing.T) {
		bigBody := strings.Repeat("x", 1024*5)

		w := makeRequest(t, app, "POST", "/api/links", bigBody)

		require.Equal(t, 413, w.Code)
	})
}

func TestApp_LongBaseURLAndShortcode(t *testing.T) {
	dbURL := setupTestDB(t)

	ctx := t.Context()

	app, err := New(ctx, config.Config{
		Env: config.EnvProduction,
		App: config.App{
			BaseURL: "http://" + strings.Repeat("x", 180),
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

	var createdLink linkDTO

	step(t, "create link", func(t *testing.T) {
		originalURL := "http://test.com"
		shortcode := strings.Repeat("щ", 32)

		w := makeRequest(t, app, "POST", "/api/links", linkDTO{
			OriginalURL: originalURL,
			Shortname:   shortcode,
		})
		require.Equal(t, 201, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)

		require.Equal(t, originalURL, resp.OriginalURL)
		require.Equal(t, shortcode, resp.Shortname)

		createdLink = resp
	})

	step(t, "read link", func(t *testing.T) {
		w := makeRequest(t, app, "GET", linkURL(createdLink.ID), nil)
		require.Equal(t, 200, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)

		require.Equal(t, createdLink.Shortname, resp.Shortname)
		require.Equal(t, createdLink.OriginalURL, resp.OriginalURL)
	})

	step(t, "resolve link", func(t *testing.T) {
		totalVisitsBefore := fetchVisitsTotalCount(t, app)

		clientIP := "1.2.3.4"

		resolveLinkURL := fmt.Sprintf("/r/%s", createdLink.Shortname)
		req, err := http.NewRequest("GET", resolveLinkURL, nil)
		require.NoError(t, err)

		req.Header.Set(gin.PlatformCloudflare, clientIP)

		w := httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		require.Equal(t, 302, w.Code)
		require.Equal(t, createdLink.OriginalURL, w.Header().Get("Location"))

		totalVisitsAfter := fetchVisitsTotalCount(t, app)
		totalVisitsInc := totalVisitsAfter - totalVisitsBefore
		require.Equal(t, 1, totalVisitsInc)
	})
}

func TestApp_ResolveLinkWithLongHeaders(t *testing.T) {
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

	var createdLink linkDTO

	step(t, "create link", func(t *testing.T) {
		originalURL := "http://test.com"
		shortcode := "test"

		w := makeRequest(t, app, "POST", "/api/links", linkDTO{
			OriginalURL: originalURL,
			Shortname:   shortcode,
		})
		require.Equal(t, 201, w.Code)

		resp := bodyAsJSON[linkDTO](t, w)

		require.Equal(t, originalURL, resp.OriginalURL)
		require.Equal(t, shortcode, resp.Shortname)

		createdLink = resp
	})

	step(t, "resolve link", func(t *testing.T) {
		totalVisitsBefore := fetchVisitsTotalCount(t, app)

		clientIP := "1.2.3.4"
		clientReferer := "http://example.com/page/" + strings.Repeat("x", 1000)
		clientUserAgent := "very-long-user-agent-" + strings.Repeat("x", 1000)

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

		totalVisitsAfter := fetchVisitsTotalCount(t, app)

		totalVisitsInc := totalVisitsAfter - totalVisitsBefore
		require.Equal(t, 1, totalVisitsInc)
	})
}

func TestApp_InvalidSort(t *testing.T) {
	dbURL := setupTestDB(t)
	clearTestDB(t)

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

	step(t, "empty db", func(t *testing.T) {
		w := makeRequest(t, app, "GET", `/api/links?sort=["unknown","ASC"]`, nil)
		require.Equal(t, 400, w.Code)

		w = makeRequest(t, app, "GET", `/api/links?sort=["id","unknown"]`, nil)
		require.Equal(t, 400, w.Code)

		w = makeRequest(t, app, "GET", `/api/link_visits?sort=["unknown","ASC"]`, nil)
		require.Equal(t, 400, w.Code)

		w = makeRequest(t, app, "GET", `/api/link_visits?sort=["id","unknown"]`, nil)
		require.Equal(t, 400, w.Code)
	})

	step(t, "add link and visit", func(t *testing.T) {
		w := makeRequest(t, app, "POST", "/api/links", linkDTO{
			OriginalURL: "http://test.com",
			Shortname:   "test",
		})
		require.Equal(t, 201, w.Code)

		clientIP := "1.2.3.4"

		req, err := http.NewRequest("GET", "/r/test", nil)
		require.NoError(t, err)

		req.Header.Set(gin.PlatformCloudflare, clientIP)

		w = httptest.NewRecorder()
		app.Router.ServeHTTP(w, req)

		require.Equal(t, 302, w.Code)
	})

	step(t, "db with data", func(t *testing.T) {
		w := makeRequest(t, app, "GET", `/api/links?sort=["unknown","ASC"]`, nil)
		require.Equal(t, 400, w.Code)

		w = makeRequest(t, app, "GET", `/api/links?sort=["id","unknown"]`, nil)
		require.Equal(t, 400, w.Code)

		w = makeRequest(t, app, "GET", `/api/link_visits?sort=["unknown","ASC"]`, nil)
		require.Equal(t, 400, w.Code)

		w = makeRequest(t, app, "GET", `/api/link_visits?sort=["id","unknown"]`, nil)
		require.Equal(t, 400, w.Code)
	})
}

func TestApp_RunAndStop(t *testing.T) {
	dbURL := setupTestDB(t)

	ctx := t.Context()

	var wg sync.WaitGroup

	addr := pickFreeAddr(t)

	app, err := New(ctx, config.Config{
		Env: config.EnvProduction,
		App: config.App{
			BaseURL: "http://example.com",
		},
		Database: config.Database{
			URL: dbURL,
		},
		HTTP: config.HTTP{
			Addr: addr,
		},
	}, discardLogger())
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	reqStarted := make(chan struct{})
	releaseReq := make(chan struct{})
	appErr := make(chan error)

	app.Router.GET("/test/wait", func(c *gin.Context) {
		close(reqStarted)

		<-releaseReq
		c.String(200, "done")
	})

	runCtx, stopApp := context.WithCancel(ctx)

	wg.Go(func() {
		appErr <- app.Run(runCtx)
	})

	waitForServer(t, addr)

	wg.Go(func() {
		resp, err := http.Get("http://" + addr + "/test/wait")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, resp.Body.Close())
		}()

		require.Equal(t, 200, resp.StatusCode)

		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, "done", string(data))
	})

	<-reqStarted
	stopApp()

	wg.Go(func() {
		// NOTE: let server start shutdown process
		time.Sleep(time.Millisecond * 100)
		close(releaseReq)
	})

	require.NoError(t, <-appErr)

	wg.Wait()
}

func TestApp_StartupErrors(t *testing.T) {
	t.Run("failed to create short url generator", func(t *testing.T) {
		_, err := New(t.Context(), config.Config{
			Env: config.EnvProduction,
			App: config.App{
				BaseURL: "invalid",
			},
			HTTP: config.HTTP{
				Addr: ":8080",
			},
		}, discardLogger())
		require.ErrorIs(t, err, ErrFailedToCreateShortURLGen)
	})

	t.Run("no db url", func(t *testing.T) {
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

	t.Run("db down", func(t *testing.T) {
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

	t.Run("invalid Sentry DSN", func(t *testing.T) {
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

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
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

var (
	testPG    *testutils.TestPostgres
	testDBURL string
)

func setupTestDB(t *testing.T) string {
	t.Helper()

	if testDBURL != "" {
		return testDBURL
	}

	ctx := t.Context()

	pg, err := testutils.SetupTestPostgres(ctx)
	require.NoError(t, err)

	require.NoError(t, pg.Migrate(migrationsDir(t)))

	testPG = &pg
	testDBURL = pg.URL

	return pg.URL
}

func clearTestDB(t *testing.T) {
	t.Helper()

	require.NotNil(t, testPG)

	err := testPG.TruncateTables(t.Context(), "visits", "links")
	require.NoError(t, err)
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

func parseTotalFromContentRange(t *testing.T, w *httptest.ResponseRecorder) int {
	t.Helper()

	contentRange := w.Header().Get("Content-Range")

	if len(contentRange) < 3 {
		t.Fatalf("Content-Range is too short: '%s'", contentRange)
	}

	for i := len(contentRange) - 2; i >= 0; i-- {
		if contentRange[i] == '/' {
			count, err := strconv.ParseInt(contentRange[i+1:], 10, 32)
			if err != nil {
				t.Fatalf("failed to parse total count in Content-Range: %v", err)
			}

			return int(count)
		}
	}

	t.Fatal("total count not found in Content-Range")

	return 0
}

func fetchVisitsTotalCount(t *testing.T, app *App) int {
	w := makeRequest(t, app, "GET", "/api/link_visits", nil)
	require.Equal(t, 200, w.Code)

	return parseTotalFromContentRange(t, w)
}

func waitForServer(t *testing.T, addr string) {
	t.Helper()

	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
		if err != nil {
			return false
		}

		require.NoError(t, conn.Close())

		return true
	}, 2*time.Second, 10*time.Millisecond)
}

func pickFreeAddr(t *testing.T) string {
	t.Helper()

	conn, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)

	addr := conn.Addr().String()

	require.NoError(t, conn.Close())

	return addr
}
