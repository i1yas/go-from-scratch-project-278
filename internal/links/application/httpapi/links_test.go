package httpapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links"
	"hexleturlshort/internal/links/application"
)

func TestGetLinkByID(t *testing.T) {
	cases := []struct {
		name       string
		url        string
		setup      func(svg *fakeLinkService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "ok",
			url:  "/api/links/101",
			setup: func(svc *fakeLinkService) {
				svc.
					On("GetLinkByID", mock.Anything, int64(101)).
					Return(
						createValidLink(t,
							101,
							"http://domain101.com",
							"link-101",
						),
						nil,
					)
			},
			wantStatus: 200,
			wantBody: `{
				"id": 101,
				"original_url": "http://domain101.com",
				"short_url": "http://short/r/link-101",
				"short_name": "link-101"
			}`,
		},
		{
			name: "not found",
			url:  "/api/links/999",
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinkByID", mock.Anything, int64(999)).
					Return(links.Link{}, application.ErrLinkNotFound)
			},
			wantStatus: 404,
			wantBody:   `{"error":"link not found"}`,
		},
		{
			name:       "invalid id",
			url:        "/api/links/test",
			setup:      func(*fakeLinkService) {},
			wantStatus: 400,
			wantBody:   `{"error":"invalid id"}`,
		},
		{
			name: "internal store error",
			url:  "/api/links/101",
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinkByID", mock.Anything, int64(101)).
					Return(links.Link{}, application.ErrStoreInternal)
			},
			wantStatus: 500,
		},
		{
			name: "invalid store value",
			url:  "/api/links/101",
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinkByID", mock.Anything, int64(101)).
					Return(links.Link{}, application.ErrInvalidStoreValue)
			},
			wantStatus: 500,
		},
		{
			name: "unknown error",
			url:  "/api/links/101",
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinkByID", mock.Anything, int64(101)).
					Return(links.Link{}, errors.New("unknown"))
			},
			wantStatus: 500,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			linksSvc := new(fakeLinkService)
			router := setupTestRouter(t, linksSvc)

			tc.setup(linksSvc)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tc.url, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if tc.wantBody != "" {
				require.JSONEq(t, tc.wantBody, w.Body.String())
			}

			linksSvc.AssertExpectations(t)
		})
	}
}

func TestGetLinks(t *testing.T) {
	cases := []struct {
		name             string
		url              string
		setup            func(svg *fakeLinkService)
		wantStatus       int
		wantBody         string
		wantContentRange string
	}{
		{
			name: "no query params",
			url:  "/api/links",
			setup: func(svc *fakeLinkService) {
				svc.
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: createRange(t, 0, 5),
						Sort:  createSort(t, "id", "ASC"),
					}).
					Return(application.LinksResult{
						Items: createLinks(t, []int64{
							1, 2, 3, 4, 5,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_links-no-params"),
			wantContentRange: "links 0-5/10",
		},
		{
			name: "links from 2nd to 3rd, got 2",
			url:  "/api/links?range=[1,3]",
			setup: func(svc *fakeLinkService) {
				svc.
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: createRange(t, 1, 3),
						Sort:  createSort(t, "id", "ASC"),
					}).
					Return(application.LinksResult{
						Items: createLinks(t, []int64{
							2, 3,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_links?range=[1,3]"),
			wantContentRange: "links 1-3/10",
		},
		{
			name: "sort by original url DESC",
			url:  `/api/links?sort=["original_url","DESC"]`,
			setup: func(svc *fakeLinkService) {
				svc.
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: createRange(t, 0, 5),
						Sort:  createSort(t, "original_url", "DESC"),
					}).
					Return(application.LinksResult{
						Items: createLinks(t, []int64{
							5, 4, 3, 2, 1,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_links?sort=[original_url,DESC]"),
			wantContentRange: "links 0-5/10",
		},
		{
			name: "sort by shot_url, link from 5th to 7th",
			url:  `/api/links?range=[5,8]&sort=["short_url","ASC"]`,
			setup: func(svc *fakeLinkService) {
				svc.
					On("GetLinks", mock.Anything, application.GetLinksParams{
						Range: createRange(t, 5, 8),
						Sort:  createSort(t, "short_url", "ASC"),
					}).
					Return(application.LinksResult{
						Items: createLinks(t, []int64{
							5, 6, 7,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_links?range=[5,8]&sort=[short_url,ASC]"),
			wantContentRange: "links 5-8/10",
		},
		{
			name:       "invalid range",
			url:        `/api/links?range=5,8`,
			setup:      func(*fakeLinkService) {},
			wantStatus: 400,
			wantBody:   `{"error":"invalid query"}`,
		},
		{
			name:       "invalid sort",
			url:        `/api/links?sort=id`,
			setup:      func(*fakeLinkService) {},
			wantStatus: 400,
			wantBody:   `{"error":"invalid query"}`,
		},
		{
			name: "internal store error",
			url:  `/api/links`,
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinks", mock.Anything, mock.Anything).
					Return(application.LinksResult{},
						application.ErrStoreInternal)
			},
			wantStatus: 500,
		},
		{
			name: "invalid store value",
			url:  `/api/links`,
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinks", mock.Anything, mock.Anything).
					Return(application.LinksResult{},
						application.ErrInvalidStoreValue)
			},
			wantStatus: 500,
		},
		{
			name: "invalid store value",
			url:  `/api/links`,
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinks", mock.Anything, mock.Anything).
					Return(application.LinksResult{},
						application.ErrInvalidStoreValue)
			},
			wantStatus: 500,
		},
		{
			name: "unknown error",
			url:  "/api/links/101",
			setup: func(svc *fakeLinkService) {
				svc.On("GetLinkByID", mock.Anything, int64(101)).
					Return(links.Link{}, errors.New("unknown"))
			},
			wantStatus: 500,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			linksSvc := new(fakeLinkService)
			router := setupTestRouter(t, linksSvc)

			tc.setup(linksSvc)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tc.url, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			require.Equal(t, tc.wantContentRange, w.Header().Get("Content-Range"))

			if tc.wantBody != "" {
				require.JSONEq(t, tc.wantBody, w.Body.String())
			}

			linksSvc.AssertExpectations(t)
		})
	}
}

func TestCreateLink(t *testing.T) {
	code := "link-101"

	cases := []struct {
		name        string
		requestBody string
		setup       func(svg *fakeLinkService)
		wantStatus  int
		wantBody    string
	}{
		{
			name: "create with shortcode",
			requestBody: `{
				"original_url": "http://domain101.com",
				"short_name": "link-101"
			}`,
			setup: func(svc *fakeLinkService) {
				svc.
					On("CreateLink", mock.Anything, application.CreateLinkParams{
						OriginalURL: "http://domain101.com",
						ShortCode:   &code,
					}).
					Return(
						createValidLink(t,
							101,
							"http://domain101.com",
							"link-101",
						),
						nil,
					)
			},
			wantStatus: 200,
			// TODO: change return value according to spec
			wantBody: `{
				"id": 101,
				"original_url": "http://domain101.com",
				"short_url": "http://short/r/link-101",
				"short_name": "link-101"
			}`,
		},
		{
			name: "create without shortcode, generated",
			requestBody: `{
				"original_url": "http://domain101.com"
			}`,
			setup: func(svc *fakeLinkService) {
				svc.
					On("CreateLink", mock.Anything, application.CreateLinkParams{
						OriginalURL: "http://domain101.com",
					}).
					Return(
						createValidLink(t,
							101,
							"http://domain101.com",
							"generated-code",
						),
						nil,
					)
			},
			wantStatus: 200,
			wantBody: `{
				"id": 101,
				"original_url": "http://domain101.com",
				"short_url": "http://short/r/generated-code",
				"short_name": "generated-code"
			}`,
		},
		{
			name: "error invalid original url",
			requestBody: `{
				"original_url": "bla"
			}`,
			setup: func(svc *fakeLinkService) {
				svc.
					On("CreateLink", mock.Anything, application.CreateLinkParams{
						OriginalURL: "bla",
					}).
					Return(links.Link{}, &links.LinkError{
						Fields: map[string]string{
							"original_url": "invalid url",
						},
					})
			},
			wantStatus: 422,
			wantBody: `{
				"errors": {
					"original_url": "invalid url"
				}
			}`,
		},
		{
			name: "shortcode conflict",
			requestBody: `{
				"original_url": "http://test.com",
				"short_name": "link-101"
			}`,
			setup: func(svc *fakeLinkService) {
				svc.
					On("CreateLink", mock.Anything, application.CreateLinkParams{
						OriginalURL: "http://test.com",
						ShortCode:   &code,
					}).
					Return(links.Link{}, application.ErrShortCodeConflict)
			},
			wantStatus: 409,
			wantBody: `{
				"errors": {
					"short_name": "link with same short_name already exist"
				}
			}`,
		},
		{
			name:        "error invalid json",
			requestBody: `}{`,
			setup:       func(_ *fakeLinkService) {},
			wantStatus:  400,
			wantBody:    `{"error": "invalid json"}`,
		},
		{
			name: "internal error",
			requestBody: `{
				"original_url": "http://domain101.com"
			}`,
			setup: func(svc *fakeLinkService) {
				svc.On("CreateLink", mock.Anything, mock.Anything).
					Return(links.Link{}, application.ErrStoreInternal)
			},
			wantStatus: 500,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			linksSvc := new(fakeLinkService)
			router := setupTestRouter(t, linksSvc)

			tc.setup(linksSvc)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/links", strings.NewReader(tc.requestBody))
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)

			if tc.wantBody != "" {
				require.JSONEq(t, tc.wantBody, w.Body.String())
			}

			linksSvc.AssertExpectations(t)
		})
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
