package service

import (
	"context"
	"errors"
	"testing"

	repoMocks "github.com/nguyendong2003/bookmark-management/internal/repository/mocks"
	serviceMocks "github.com/nguyendong2003/bookmark-management/internal/service/mocks"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLService_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRepo    func(ctx context.Context) *repoMocks.URLStorage
		setupCodeGen func(ctx context.Context) *serviceMocks.Password

		expectedCode string
		expectedErr  error
	}{
		{
			name: "normal case - new key",

			setupRepo: func(ctx context.Context) *repoMocks.URLStorage {
				mock := repoMocks.NewURLStorage(t)

				mock.On("StoreURL", ctx, "abc1234567", "https://google.com").Return(nil)

				return mock
			},

			setupCodeGen: func(ctx context.Context) *serviceMocks.Password {
				mock := serviceMocks.NewPassword(t)

				mock.On("GeneratePassword").Return("abc1234567", nil)

				return mock
			},

			expectedCode: "abc1234567",
			expectedErr:  nil,
		},
		{
			name: "password generate error",

			setupRepo: func(ctx context.Context) *repoMocks.URLStorage {
				mock := repoMocks.NewURLStorage(t)
				return mock
			},

			setupCodeGen: func(ctx context.Context) *serviceMocks.Password {
				mock := serviceMocks.NewPassword(t)

				mock.On("GeneratePassword").Return("", errors.New("generate error"))

				return mock
			},

			expectedCode: "",
			expectedErr:  errors.New("generate error"),
		},
		{
			name: "repository store error",

			setupRepo: func(ctx context.Context) *repoMocks.URLStorage {
				mock := repoMocks.NewURLStorage(t)

				mock.On("StoreURL", ctx, "abc1234567", "https://google.com").Return(errors.New("redis error"))

				return mock
			},

			setupCodeGen: func(ctx context.Context) *serviceMocks.Password {
				mock := serviceMocks.NewPassword(t)
				mock.On("GeneratePassword").Return("abc1234567", nil)
				return mock
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

			mockRepo := tc.setupRepo(ctx)
			mockCodeGen := tc.setupCodeGen(ctx)

			service := NewShortenURL(mockRepo, mockCodeGen)

			code, err := service.ShortenURL(ctx, "https://google.com")

			assert.Equal(t, tc.expectedCode, code)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockCodeGen.AssertExpectations(t)
		})
	}
}

func TestShortenURLService_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRepo func(ctx context.Context) *repoMocks.URLStorage

		expectedURL string
		expectedErr error
	}{
		{
			name: "success",

			setupRepo: func(ctx context.Context) *repoMocks.URLStorage {
				mock := repoMocks.NewURLStorage(t)

				mock.On("GetURL", ctx, "abc1234567").Return("https://google.com", nil)

				return mock
			},

			expectedURL: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "code not exist",

			setupRepo: func(ctx context.Context) *repoMocks.URLStorage {
				mock := repoMocks.NewURLStorage(t)

				mock.On("GetURL", ctx, "abc1234567").Return("", redis.Nil)

				return mock
			},

			expectedURL: "",
			expectedErr: ErrCodeNotExist,
		},
		{
			name: "repository error",

			setupRepo: func(ctx context.Context) *repoMocks.URLStorage {
				mock := repoMocks.NewURLStorage(t)

				mock.On("GetURL", ctx, "abc1234567").Return("", errors.New("redis error"))

				return mock
			},

			expectedURL: "",
			expectedErr: errors.New("redis error"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			mockRepo := tc.setupRepo(ctx)

			service := NewShortenURL(mockRepo, nil)

			url, err := service.GetURL(ctx, "abc1234567")

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
