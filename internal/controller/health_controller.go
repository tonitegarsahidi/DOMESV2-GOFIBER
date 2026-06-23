package controller

import (
	"github.com/gofiber/fiber/v2"
	"domesv2/internal/service"
	"domesv2/pkg/response"
)

type HealthController struct {
	healthService service.HealthService
}

func NewHealthController(healthService service.HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}

func (ctrl *HealthController) Check(c *fiber.Ctx) error {
	health, err := ctrl.healthService.CheckHealth()
	if err != nil {
		return response.Error(c, err)
	}

	if health.Status != "healthy" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(health)
	}

	return response.Success(c, health, "All systems operational")
}
