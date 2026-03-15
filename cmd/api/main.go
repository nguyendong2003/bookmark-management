package main

import (
	"github.com/nguyendong2003/bookmark-management/internal/api"
	"github.com/nguyendong2003/bookmark-management/pkg/logger"
	redisPkg "github.com/nguyendong2003/bookmark-management/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0
// @description This is the API documentation for the Bookmark Management Service.
// @host localhost:8080
// @BasePath /
func main() {
	logger.SetLogLevel()

	cfg, err := api.NewConfig("BOOKMARK_SERVICE")
	if err != nil {
		panic(err)
	}

	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(cfg, redisClient)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
