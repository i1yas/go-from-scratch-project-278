package httpapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"hexleturlshort/internal/links/application"
	"hexleturlshort/internal/links/application/testutils"
)

func TestGetVisits(t *testing.T) {
	cases := []struct {
		name             string
		url              string
		setup            func(svg *fakeVisitsService)
		wantStatus       int
		wantBody         string
		wantContentRange string
	}{
		{
			name: "no query params, default params used",
			url:  "/api/link_visits",
			setup: func(svc *fakeVisitsService) {
				svc.
					On("GetVisits", mock.Anything, application.GetVisitsParams{
						Range: testutils.Range(t, 0, 4),
						Sort:  testutils.SortOrder(t, "id", "ASC"),
					}).
					Return(application.VisitsResult{
						Items: createVisits(t, []int64{
							1, 2, 3, 4, 5,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_visits-no-params"),
			wantContentRange: "visits 0-4/10",
		},
		{
			name: "visits from 2nd to 3rd got 3",
			url:  "/api/link_visits?range=[1,3]",
			setup: func(svc *fakeVisitsService) {
				svc.
					On("GetVisits", mock.Anything, application.GetVisitsParams{
						Range: testutils.Range(t, 1, 3),
						Sort:  testutils.SortOrder(t, "id", "ASC"),
					}).
					Return(application.VisitsResult{
						Items: createVisits(t, []int64{
							2, 3, 4,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_visits?range=[1,3]"),
			wantContentRange: "visits 1-3/10",
		},
		{
			name: "sort by IP DESC",
			url:  `/api/link_visits?sort=["ip", "DESC"]`,
			setup: func(svc *fakeVisitsService) {
				svc.
					On("GetVisits", mock.Anything, application.GetVisitsParams{
						Range: testutils.Range(t, 0, 4),
						Sort:  testutils.SortOrder(t, "ip", "DESC"),
					}).
					Return(application.VisitsResult{
						Items: createVisits(t, []int64{
							10, 1, 2, 3, 4,
						}),
						Total: 10,
					}, nil)
			},
			wantStatus:       200,
			wantBody:         loadResponseFixture(t, "get_visits?sort=[ip,DESC]"),
			wantContentRange: "visits 0-4/10",
		},
		{
			name:       "invalid range",
			url:        "/api/link_visits?range=5,8",
			setup:      func(*fakeVisitsService) {},
			wantStatus: 400,
			wantBody:   `{"error":"invalid query"}`,
		},
		{
			name:       "invalid sort",
			url:        `/api/link_visits?sort=id`,
			setup:      func(*fakeVisitsService) {},
			wantStatus: 400,
			wantBody:   `{"error":"invalid query"}`,
		},
		{
			name: "internal store error",
			url:  `/api/link_visits`,
			setup: func(svc *fakeVisitsService) {
				svc.On("GetVisits", mock.Anything, mock.Anything).
					Return(application.VisitsResult{},
						application.ErrStoreInternal)
			},
			wantStatus: 500,
		},
		{
			name: "invalid store value",
			url:  `/api/link_visits`,
			setup: func(svc *fakeVisitsService) {
				svc.On("GetVisits", mock.Anything, mock.Anything).
					Return(application.VisitsResult{},
						application.ErrInvalidStoreValue)
			},
			wantStatus: 500,
		},
		{
			name: "unknown error",
			url:  "/api/link_visits",
			setup: func(svc *fakeVisitsService) {
				svc.On("GetVisits", mock.Anything, mock.Anything).
					Return(application.VisitsResult{},
						errors.New("unknown"))
			},
			wantStatus: 500,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			visitsSvc := new(fakeVisitsService)
			router := setupTestRouter(t, &fakeLinkService{}, visitsSvc)

			tc.setup(visitsSvc)

			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tc.url, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			require.Equal(t, tc.wantContentRange, w.Header().Get("Content-Range"))

			if tc.wantBody != "" {
				require.JSONEq(t, tc.wantBody, w.Body.String())
			}

			visitsSvc.AssertExpectations(t)
		})
	}
}
