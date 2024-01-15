package cache

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/tests"
)

func TestCacheClient(t *testing.T) {
	type testData struct {
		Value string `json:"value"`
	}
	
	c := CreateClient(context.Background(), cacheAdapter.Redis, nil, tests.CreateRedisConnection(t))
	
	t.Run(
		"set", func(t *testing.T) {
			assert.Nil(t, c.Set("test", testData{Value: "test"}, time.Minute))
		},
	)
	
	t.Run(
		"get", func(t *testing.T) {
			var d testData
			assert.Nil(t, c.Get("test", &d))
			assert.Equal(t, "test", d.Value)
		},
	)
}
