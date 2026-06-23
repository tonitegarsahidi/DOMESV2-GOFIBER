package service

import (
	"time"

	"domesv2/config"
	"domesv2/internal/model"
	"domesv2/internal/repository"
	"go.uber.org/zap"
)

type HealthService interface {
	CheckHealth() (*model.HealthCheckResponse, error)
}

type healthService struct {
	healthRepo repository.HealthRepository
	cfg        *config.Config
}

func NewHealthService(healthRepo repository.HealthRepository) HealthService {
	return &healthService{
		healthRepo: healthRepo,
		cfg:        config.AppConfig,
	}
}

func (s *healthService) CheckHealth() (*model.HealthCheckResponse, error) {
	services := make(map[string]string)
	overallStatus := "healthy"

	// Check database
	dbErr := s.healthRepo.CheckDatabase()
	if dbErr != nil {
		services["database"] = "unhealthy"
		services["database_error"] = dbErr.Error()
		overallStatus = "unhealthy"
		zap.L().Error("Database health check failed", zap.Error(dbErr))
	} else {
		services["database"] = "healthy"
	}

	// Check Redis if enabled
	if s.cfg.Redis.Enabled {
		redisErr := s.healthRepo.CheckRedis()
		if redisErr != nil {
			services["redis"] = "unhealthy"
			services["redis_error"] = redisErr.Error()
			overallStatus = "unhealthy"
			zap.L().Error("Redis health check failed", zap.Error(redisErr))
		} else {
			services["redis"] = "healthy"
		}
	} else {
		services["redis"] = "disabled"
	}

	// Application status
	services["application"] = "healthy"

	return &model.HealthCheckResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  services,
	}, nil
}
