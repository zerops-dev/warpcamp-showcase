package redis

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zerops-dev/warpcamp-showcase/src/env"
)

const (
	EnvHost = "REDIS_HOST"
	EnvPort = "REDIS_PORT"
)

func NewRedis() (*redis.Client, error) {
	if err := env.Check(EnvHost, EnvPort); err != nil {
		return nil, err
	}

	port, err := strconv.ParseUint(os.Getenv(EnvPort), 10, 16)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(&redis.Options{
		Network:         "tcp",
		Addr:            fmt.Sprintf("%s:%d", os.Getenv(EnvHost), port),
		DB:              1,
		MaxRetries:      5,
		MinRetryBackoff: 5,
		MaxRetryBackoff: 5,
		DialTimeout:     time.Second * 2,
		ReadTimeout:     time.Second * 2,
		WriteTimeout:    time.Second * 2,
	}), nil
}
