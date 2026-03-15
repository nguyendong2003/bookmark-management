package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenURLService_ShortenURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLen int
		expectedErr error
	}{
		{
			name: "normal case",

			expectedLen: 10,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testService := NewPassword()

			password, err := testService.GeneratePassword()

			assert.Equal(t, tc.expectedLen, len(password))
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, urlSafeRegex.MatchString(password), true)

		})
	}
}
