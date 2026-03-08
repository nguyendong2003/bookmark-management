package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func() *mocks.HealthCheck

		expectedStatus   int
		expectedResponse map[string]any
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setupMockService: func() *mocks.HealthCheck {
				serviceMock := mocks.NewHealthCheck(t)

				serviceMock.On("CheckHealth").Return(map[string]any{
					"message":      "OK",
					"service_name": "test-service",
					"instance_id":  "test-instance-id",
				})
				return serviceMock
			},

			expectedStatus: http.StatusOK,
			expectedResponse: map[string]any{
				"message":      "OK",
				"service_name": "test-service",
				"instance_id":  "test-instance-id",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)

			mockService := tc.setupMockService()
			testHandler := NewHealthCheck(mockService)

			testHandler.CheckHealth(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			// convert expected response to JSON
			expectedJSON, err := json.Marshal(tc.expectedResponse)
			assert.NoError(t, err)

			// assert body
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}

}
