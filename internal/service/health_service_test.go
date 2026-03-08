package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckService_CheckHealth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		serviceName string
		instanceID  string

		expected map[string]any
	}{
		{
			name: "normal case",

			serviceName: "test-service",
			instanceID:  "instance-1",

			expected: map[string]any{
				"message":      "OK",
				"service_name": "test-service",
				"instance_id":  "instance-1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testService := NewHealthCheck(tc.serviceName, tc.instanceID)

			result := testService.CheckHealth()

			assert.Equal(t, tc.expected, result)
		})
	}
}
