package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"market_system/api"
	"market_system/config"
	"market_system/storage/postgres"
	"market_system/storage/redis"
)

func main() {

	var cfg = config.Load()

	pgStorage, err := postgres.NewConnectionPostgres(&cfg)
	if err != nil {
		panic(err)
	}

	cache, err := redis.NewConnectionRedis(&cfg)
	if err != nil {
		panic(err)
	}

	// gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	api.SetUpApi(r, &cfg, pgStorage, cache)

	log.Println("Listening:", cfg.ServiceHost+cfg.ServiceHTTPPort, "...")
	if err := r.Run(cfg.ServiceHost + cfg.ServiceHTTPPort); err != nil {
		panic("Listent and service panic:" + err.Error())
	}
}
