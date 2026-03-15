package service

import (
	"context"
	"errors"
	"testing"

	repoMocks "github.com/nguyendong2003/bookmark-management/internal/repository/mocks"
	serviceMocks "github.com/nguyendong2003/bookmark-management/internal/service/mocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURLService_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(repo *repoMocks.URLStorage, pass *serviceMocks.Password)

		expectedCode string
		expectedErr  error
	}{
		{
			name: "success",

			setupMock: func(repo *repoMocks.URLStorage, pass *serviceMocks.Password) {
				pass.On("GeneratePassword").
					Return("abc1234567", nil)

				repo.On("StoreURL", mock.Anything, "abc1234567", "https://google.com").
					Return(nil)
			},

			expectedCode: "abc1234567",
			expectedErr:  nil,
		},
		{
			name: "password generate error",

			setupMock: func(repo *repoMocks.URLStorage, pass *serviceMocks.Password) {
				pass.On("GeneratePassword").
					Return("", errors.New("generate error"))
			},

			expectedCode: "",
			expectedErr:  errors.New("generate error"),
		},
		{
			name: "repository store error",

			setupMock: func(repo *repoMocks.URLStorage, pass *serviceMocks.Password) {
				pass.On("GeneratePassword").
					Return("abc1234567", nil)

				repo.On("StoreURL", mock.Anything, "abc1234567", "https://google.com").
					Return(errors.New("redis error"))
			},

			expectedCode: "",
			expectedErr:  errors.New("redis error"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockRepo := repoMocks.NewURLStorage(t)
			mockPassword := serviceMocks.NewPassword(t)

			tc.setupMock(mockRepo, mockPassword)

			service := NewShortenURL(mockRepo, mockPassword)

			code, err := service.ShortenURL(ctx, "https://google.com")

			assert.Equal(t, tc.expectedCode, code)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockPassword.AssertExpectations(t)
		})
	}
}

func TestShortenURLService_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(repo *repoMocks.URLStorage)

		expectedURL string
		expectedErr error
	}{
		{
			name: "success",

			setupMock: func(repo *repoMocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "abc").
					Return("https://google.com", nil)
			},

			expectedURL: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "code not exist",

			setupMock: func(repo *repoMocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "abc").
					Return("", redis.Nil)
			},

			expectedURL: "",
			expectedErr: ErrCodeNotExist,
		},
		{
			name: "repository error",

			setupMock: func(repo *repoMocks.URLStorage) {
				repo.On("GetURL", mock.Anything, "abc").
					Return("", errors.New("redis error"))
			},

			expectedURL: "",
			expectedErr: errors.New("redis error"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockRepo := repoMocks.NewURLStorage(t)

			tc.setupMock(mockRepo)

			service := NewShortenURL(mockRepo, nil)

			url, err := service.GetURL(ctx, "abc")

			assert.Equal(t, tc.expectedURL, url)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
