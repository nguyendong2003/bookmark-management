package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

func TestShortenURLHandler_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShortenURL

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case - success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://www.youtube.com/"}`)),
				)
			},

			setupMockService: func(ctx context.Context) *mocks.ShortenURL {
				serviceMock := mocks.NewShortenURL(t)
				serviceMock.On("ShortenURL", ctx, "https://www.youtube.com/").Return("123456", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: `{"key":"123456"}`,
		},
		{
			name: "fail case - service failed",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBuffer([]byte(`{"url":"https://www.youtube.com/"}`)),
				)
			},

			setupMockService: func(ctx context.Context) *mocks.ShortenURL {
				serviceMock := mocks.NewShortenURL(t)
				serviceMock.On("ShortenURL", ctx, "https://www.youtube.com/").Return("", testErr)
				return serviceMock
			},

			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"error":"Internal server error"}`,
		},
		{
			name: "fail case - bad input",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodPost,
					"/link/shorten",
					bytes.NewBuffer([]byte(`{"url":"not a link"}`)),
				)
			},

			setupMockService: func(ctx context.Context) *mocks.ShortenURL {
				serviceMock := mocks.NewShortenURL(t)
				return serviceMock
			},

			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid input"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)

			mockService := tc.setupMockService(gc)
			testHandler := NewShortenURL(mockService)

			testHandler.ShortenURL(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
			// assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}
