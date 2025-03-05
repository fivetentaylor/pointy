package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, fmt.Errorf("error getting redis url: %w", err)
	}
	rc := redis.NewClient(opt)
	_, err = rc.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("error pinging redis: %w", err)
	}

	return rc, nil
}
