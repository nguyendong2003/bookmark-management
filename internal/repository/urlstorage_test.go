package repository

import (
	"context"
	"testing"

	redisPkg "github.com/nguyendong2003/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func() *redis.Client
		inputCode string
		inputURL  string

		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client, inputCode, inputURL string)
	}{
		{
			name: "normal case",

			setupMock: func() *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			inputCode: "123",
			inputURL:  "https://google.com",

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client, inputCode, inputURL string) {
				res, err := r.Get(ctx, inputCode).Result()
				assert.NoError(t, err)
				assert.Equal(t, inputURL, res)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisClient := tc.setupMock()
			testRepo := NewURLStorage(redisClient)

			err := testRepo.StoreURL(ctx, tc.inputCode, tc.inputURL)
			if err == nil {
				tc.verifyFunc(ctx, redisClient, tc.inputCode, tc.inputURL)
			}

		})
	}
}

func TestURLStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client
		inputCode string

		expectedURL string
		expectedErr error
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)

				err := mock.Set(ctx, "123", "https://google.com", 0).Err()
				assert.NoError(t, err)

				return mock
			},

			inputCode:   "123",
			expectedURL: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "fail case - key not exist",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},

			inputCode:   "999",
			expectedURL: "",
			expectedErr: redis.Nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			redisClient := tc.setupMock(ctx)
			testRepo := NewURLStorage(redisClient)

			url, err := testRepo.GetURL(ctx, tc.inputCode)

			assert.Equal(t, tc.expectedURL, url)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
