package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHandler_GenPass(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func() *mocks.Password

		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},

			setupMockService: func() *mocks.Password {
				serviceMock := mocks.NewPassword(t)

				// Viết mock cho phương thức "GenPass" của interface Password trong package handler.
				// Vì phương thức này gọi phương thức "GeneratePassword" của interface Password trong package service, nên cần tạo mock service
				// "GeneratePassword" là tên của method của interface Password trong package service (nằm ở file password_service.go),
				// và "Return" là phương thức của mock để định nghĩa giá trị trả về khi method đó được gọi.
				// Hàm "GeneratePassword" của interface Password trong package service trả về (string, error), nên chúng ta cần cung cấp hai giá trị trong "Return": một chuỗi đại diện cho mật khẩu được tạo ra và một giá trị lỗi (ở đây là nil để biểu thị không có lỗi).
				serviceMock.On("GeneratePassword").Return("123456789", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: "123456789",
		},
		{
			name: "fail",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},

			setupMockService: func() *mocks.Password {
				serviceMock := mocks.NewPassword(t)

				// Viết mock cho phương thức "GenPass" của interface Password trong package handler.
				// Vì phương thức này gọi phương thức "GeneratePassword" của interface Password trong package service, nên cần tạo mock service
				// "GeneratePassword" là tên của method của interface Password trong package service (nằm ở file password_service.go),
				// và "Return" là phương thức của mock để định nghĩa giá trị trả về khi method đó được gọi.
				// Hàm "GeneratePassword" của interface Password trong package service trả về (string, error), nên chúng ta cần cung cấp hai giá trị trong "Return": một chuỗi đại diện cho mật khẩu được tạo ra và một giá trị lỗi (ở đây là nil để biểu thị không có lỗi).
				serviceMock.On("GeneratePassword").Return("123456789", nil)
				return serviceMock
			},

			expectedStatus:   http.StatusOK,
			expectedResponse: "12345678",
		},
		{
			name: "internal server error",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},

			setupMockService: func() *mocks.Password {
				serviceMock := mocks.NewPassword(t)
				serviceMock.On("GeneratePassword").Return("", errors.New("Failed Generate Password"))
				return serviceMock
			},

			expectedStatus: http.StatusInternalServerError,

			// Giả sử serviceMock trả về error: "Failed Generate Password", nhưng expectedResponse: "Failed" nên testcase này fail
			expectedResponse: "Failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// rec là viết tắt của recorder. gc là viết tắt của gin context
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)

			mockService := tc.setupMockService()
			testHandler := NewPassword(mockService)

			testHandler.GenPass(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}
