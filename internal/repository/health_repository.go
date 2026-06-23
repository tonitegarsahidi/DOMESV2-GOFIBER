package repository

import (
	"fmt"

	"domesv2/config/database"
	"domesv2/config/redis"
)

type HealthRepository interface {
	CheckDatabase() error
	CheckRedis() error
}

type healthRepository struct{}

func NewHealthRepository() HealthRepository {
	return &healthRepository{}
}

func (r *healthRepository) CheckDatabase() error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (r *healthRepository) CheckRedis() error {
	if !redis.IsRedisEnabled() {
		return nil
	}

	ctx := redis.GetCtx()
	_, err := redis.GetRedis().Ping(ctx).Result()
	return err
}
