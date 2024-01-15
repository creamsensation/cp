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

func TestFlashMessage(t *testing.T) {
	c := &control{
		context: context.Background(),
		core: &core{
			redis: tests.CreateRedisConnection(t),
		},
		request:  httptest.NewRequest(http.MethodGet, "/test", nil),
		response: httptest.NewRecorder(),
	}
	c.core.router = createRouter(c.core)
	c.config.Cache.Adapter = cacheAdapter.Redis
	c.flash = &flashMessenger{control: c, messages: make([]FlashMessage, 0)}
	c.Flash().Add().Success("test")
	c.flash.store()
	assert.True(t, len(c.Flash().Get()) == 1)
}
