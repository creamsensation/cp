package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	
	"github.com/go-redis/redis/v8"
	
	"github.com/creamsensation/cp/internal/cache/memory"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
)

type Client interface {
	Exists(key string) bool
	Get(key string, data any) error
	Set(key string, data any, expiration time.Duration) error
	Destroy(key string) error
}

type client struct {
	ctx     context.Context
	adapter string
	memory  memory.Client
	redis   *redis.Client
}

func CreateClient(ctx context.Context, adapter string, memory memory.Client, redis *redis.Client) Client {
	return &client{
		ctx:     ctx,
		adapter: adapter,
		memory:  memory,
		redis:   redis,
	}
}

func (c client) Exists(key string) bool {
	if c.isNil() {
		return false
	}
	switch c.adapter {
	case cacheAdapter.Memory:
		return c.memory.Exists(key)
	case cacheAdapter.Redis:
		cmd := c.redis.Exists(c.ctx, key)
		if cmd == nil {
			return false
		}
		return cmd.Val() > 0
	default:
		return false
	}
}

func (c client) Get(key string, data any) error {
	if c.isNil() {
		return errors.New("adapter instance not exist")
	}
	var value string
	switch c.adapter {
	case cacheAdapter.Memory:
		value = c.memory.Get(key)
	case cacheAdapter.Redis:
		stored := c.redis.Get(c.ctx, key)
		if stored == nil {
			return nil
		}
		value = stored.Val()
	}
	if len(value) > 0 {
		return json.Unmarshal([]byte(value), data)
	}
	return nil
}

func (c client) Set(key string, data any, expiration time.Duration) error {
	if c.isNil() {
		return errors.New("adapter instance not exist")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	switch c.adapter {
	case cacheAdapter.Memory:
		return c.memory.Set(key, string(b), expiration)
	case cacheAdapter.Redis:
		if c.redis.Set(c.ctx, key, string(b), expiration).Err() != nil {
			return err
		}
		return nil
	}
	return nil
}

func (c client) Destroy(key string) error {
	switch c.adapter {
	case cacheAdapter.Memory:
		return c.memory.Destroy(key)
	case cacheAdapter.Redis:
		return c.Set(key, nil, time.Millisecond)
	}
	return nil
}

func (c client) isNil() bool {
	switch c.adapter {
	case cacheAdapter.Memory:
		return c.memory == nil
	case cacheAdapter.Redis:
		return c.redis == nil
	default:
		return true
	}
}
