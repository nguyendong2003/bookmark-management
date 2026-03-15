package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlSafeRegex = regexp.MustCompile("^[A-Za-z0-9]+$")

func TestPasswordService_GeneratePassword(t *testing.T) {
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
		tc := tc

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
