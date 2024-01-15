package cp

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/tests"
)

func TestCache(t *testing.T) {
	type testData struct {
		Value string `json:"value"`
	}
	c := &control{
		context: context.Background(),
		core:    &core{redis: tests.CreateRedisConnection(t)},
	}
	c.config.Cache.Adapter = cacheAdapter.Redis
	
	t.Run(
		"set", func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			c.Cache().Set("test", testData{Value: "test"}, time.Minute)
		},
	)
	
	t.Run(
		"get", func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			var data testData
			c.Cache().Get("test", &data)
			assert.Equal(t, "test", data.Value)
		},
	)
	
	t.Run(
		"exist", func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			assert.True(t, c.Cache().Exists("test"))
		},
	)
	
	t.Run(
		"destroy", func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			c.Cache().Destroy("test")
		},
	)
}
