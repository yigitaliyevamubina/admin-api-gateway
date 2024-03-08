package main

import (
	"fmt"
	"myproject/admin-api-gateway/api"
	"myproject/admin-api-gateway/config"
	"myproject/admin-api-gateway/pkg/db"
	"myproject/admin-api-gateway/pkg/logger"
	"myproject/admin-api-gateway/services"
	"myproject/admin-api-gateway/storage/postgres"
	"myproject/admin-api-gateway/storage/redis"

	rds "github.com/gomodule/redigo/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "admin-api-gateway")

	serviceManager, err := services.NewServiceManager(cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
		return
	}

	redisPool := rds.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (rds.Conn, error) {
			c, err := rds.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	db, _, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("cannot connect to DB", logger.Error(err))
		return
	}

	server := api.New(api.Option{
		InMemory:       redis.NewRedisRepo(&redisPool),
		Cfg:            cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		Postgres:       postgres.NewAdminRepo(db),
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		return
	}
}
