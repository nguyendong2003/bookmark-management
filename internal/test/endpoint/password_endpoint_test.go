package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyendong2003/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestPasswordEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHttp func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus  int
		expectedRespLen int
	}{
		{
			name: "success",

			setupTestHttp: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedRespLen: 10,
		},
		{
			name: "fail",

			setupTestHttp: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus:  http.StatusOK,
			expectedRespLen: 9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(&api.Config{})
			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedRespLen, rec.Body.Len())
		})
	}
}
