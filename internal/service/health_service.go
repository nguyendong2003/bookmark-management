package service

type HealthCheck interface {
	CheckHealth() map[string]any
}

type healthCheck struct {
	serviceName string
	instanceID  string
}

//go:generate go run github.com/vektra/mockery/v2@latest --name HealthCheck --filename health_service.go
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheck{
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

func (s *healthCheck) CheckHealth() map[string]any {
	return map[string]any{
		"message":      "OK",
		"service_name": s.serviceName,
		"instance_id":  s.instanceID,
	}
}
