package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewClient(addr, password, DB string) *redis.Client {
	db, err := strconv.Atoi(DB)
	if err != nil {
		log.Printf("NewCLient.db error: %v", err)
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
