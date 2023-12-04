package cp

import (
	"time"
	
	. "github.com/creamsensation/cp/internal/cache"
)

type Cache interface {
	Exists(key string) bool
	Get(key string, data any)
	Set(key string, data any, expiration time.Duration)
	Destroy(key string)
}

type cache struct {
	control *control
	client  Client
}

func createCache(control *control) Cache {
	return &cache{
		control: control,
		client: CreateClient(
			control.context,
			control.config.Cache.Adapter,
			control.core.memory,
			control.core.redis,
		),
	}
}

func (c cache) Exists(key string) bool {
	return c.client.Exists(key)
}

func (c cache) Get(key string, data any) {
	c.control.Error().Check(c.client.Get(key, data))
}

func (c cache) Set(key string, data any, expiration time.Duration) {
	c.control.Error().Check(c.client.Set(key, data, expiration))
}

func (c cache) Destroy(key string) {
	c.control.Error().Check(c.client.Destroy(key))
}
