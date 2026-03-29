package main

import (
	"context"
	"time"

	"github.com/nguyendong2003/bookmark-management/pkg/redis"
)

func main() {
	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	redisClient.Set(context.Background(), "test_key", "test_value", time.Hour)
}
