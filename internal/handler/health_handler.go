package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type HealthCheck interface {
	CheckHealth(c *gin.Context)
}

type healthCheck struct {
	healthCheckService service.HealthCheck
}

func NewHealthCheck(healthCheckService service.HealthCheck) HealthCheck {
	return &healthCheck{
		healthCheckService: healthCheckService,
	}
}

func (h *healthCheck) CheckHealth(c *gin.Context) {
	result := h.healthCheckService.CheckHealth()
	c.JSON(http.StatusOK, result)
}
