package cp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/responder/hx"
	"github.com/creamsensation/cp/internal/tests"
	"github.com/creamsensation/gox"
)

func TestHandleHx(t *testing.T) {
	c := &control{
		context: context.Background(),
		core: &core{
			redis: tests.CreateRedisConnection(t),
		},
		request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		response: httptest.NewRecorder(),
	}
	c.config.Cache.Adapter = cacheAdapter.Redis
	c.hxResponse = hx.New(c.request, c.response)
	c.core.router = createRouter(c.core)
	
	t.Run(
		"click", func(t *testing.T) {
			attrs := gox.Render(c.Handle().Hx().Click("test"))
			assert.True(t, strings.Contains(attrs, `id="hx-handler-test"`))
			assert.True(t, strings.Contains(attrs, `hx-post="/test"`))
			assert.True(t, strings.Contains(attrs, `hx-trigger="click"`))
		},
	)
	
	t.Run(
		"submit", func(t *testing.T) {
			attrs := gox.Render(c.Handle().Hx().Submit("test"))
			assert.True(t, strings.Contains(attrs, `id="hx-handler-test"`))
			assert.True(t, strings.Contains(attrs, `hx-post="/test"`))
			assert.True(t, strings.Contains(attrs, `hx-trigger="submit"`))
		},
	)
}
