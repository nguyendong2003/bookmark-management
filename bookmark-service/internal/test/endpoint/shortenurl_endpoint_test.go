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
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHttp func(api api.Engine) *httptest.ResponseRecorder
		setupRedis    func(ctx context.Context) *redis.Client

		expectedStatus int
		validateResp   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
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

			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				mock.Set(ctx, "1234567", "https://google.com", 300)
				return mock
			},

			expectedStatus: http.StatusOK,

			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				resp := map[string]string{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Len(t, resp["key"], 10, "expected key length to be 10")
			},
		},

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

			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				mock.Set(ctx, "1234567", "https://google.com", 300)
				return mock
			},

			expectedStatus: http.StatusBadRequest,

			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.JSONEq(t, `{"error":"Invalid input"}`, rec.Body.String())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mockRedis := tc.setupRedis(ctx)
			app := api.NewEngine(&api.Config{}, mockRedis)

			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.validateResp != nil {
				tc.validateResp(t, rec)
			}
		})
	}
}

func TestGetURLEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupTestHttp func(api api.Engine) *httptest.ResponseRecorder
		setupRedis    func(ctx context.Context) *redis.Client

		expectedStatus int
		validateResp   func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "fail case - code not exist",

			setupTestHttp: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(
					http.MethodGet,
					"/link/redirect/notexist",
					nil,
				)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)

				return rec
			},

			setupRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectedStatus: http.StatusBadRequest,

			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.JSONEq(t, `{"error":"Code not exist"}`, rec.Body.String())
			},
		},
		{
			name: "success",

			setupTestHttp: func(api api.Engine) *httptest.ResponseRecorder {

				// create shorten url first
				shortenReq := httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBufferString(`{"url":"https://www.youtube.com/"}`),
				)

				shortenRec := httptest.NewRecorder()
				api.ServeHTTP(shortenRec, shortenReq)

				require.Equal(t, http.StatusOK, shortenRec.Code)

				resp := map[string]string{}
				require.NoError(t, json.Unmarshal(shortenRec.Body.Bytes(), &resp))

				key := resp["key"]
				require.NotEmpty(t, key)

				// call redirect endpoint with code
				redirectReq := httptest.NewRequest(
					http.MethodGet,
					"/link/redirect/"+key,
					nil,
				)

				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, redirectReq)

				return rec
			},

			setupRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectedStatus: http.StatusMovedPermanently,

			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, "https://www.youtube.com/", rec.Header().Get("Location"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			mockRedis := tc.setupRedis(ctx)
			app := api.NewEngine(&api.Config{}, mockRedis)

			rec := tc.setupTestHttp(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.validateResp != nil {
				tc.validateResp(t, rec)
			}
		})
	}
}
