package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyendong2003/bookmark-management/internal/api"
	redisPkg "github.com/nguyendong2003/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	redisClient, err := redisPkg.NewClient("")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, redisClient.Close())
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		t.Skipf("skip integration test because redis is not available: %v", err)
	}

	testCases := []struct {
		name string

		setupTestHttp func(api api.Engine) *httptest.ResponseRecorder

		expectedStatus int
		validateResp   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "bad input - invalid url",

			setupTestHttp: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBufferString(`{"url":"not a link"}`),
				)
				// req.Header.Set("Content-Type", "application/json")

				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusBadRequest,

			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.JSONEq(t, `{"error":"Invalid input"}`, rec.Body.String())
			},
		},
		{
			name: "success",

			setupTestHttp: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBufferString(`{"url":"https://www.youtube.com/"}`),
				)
				// req.Header.Set("Content-Type", "application/json")

				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},

			expectedStatus: http.StatusOK,

			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				resp := map[string]string{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Len(t, resp["key"], 10, "expected key length to be 10")
			},
		},
	}

	app := api.NewEngine(&api.Config{}, redisClient)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.validateResp != nil {
				tc.validateResp(t, rec)
			}
		})
	}
}
