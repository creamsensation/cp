package tests

import (
	"context"
	"testing"
	
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func CreateRedisConnection(t *testing.T) *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB:   10,
		},
	)
	assert.Nil(t, client.Ping(context.Background()).Err())
	return client
}
