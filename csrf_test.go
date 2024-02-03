package cp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/tests"
)

func TestCsrf(t *testing.T) {
	var token string
	ctrl := &control{
		context: context.Background(),
		core: &core{
			form:  createFormManager(),
			redis: tests.CreateRedisConnection(t),
		},
		request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		response: httptest.NewRecorder(),
	}
	ctrl.config.Cache.Adapter = cacheAdapter.Redis
	c := createCsrf(ctrl)
	
	t.Cleanup(
		func() {
			c.Destroy(token)
		},
	)
	
	t.Run(
		"create", func(t *testing.T) {
			token = c.Create("test", "test-form", "1.1.1.1", "test")
			assert.True(t, len(token) > 0)
		},
	)
	
	t.Run(
		"get", func(t *testing.T) {
			r := c.Get(token, "test-form")
			assert.Equal(t, "1.1.1.1", r.Ip)
			assert.Equal(t, "test", r.UserAgent)
		},
	)
}
