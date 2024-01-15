package cp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	c := &control{
		request: httptest.NewRequest(http.MethodGet, "/test?a=test&b=1&c=true&d=1.99", nil),
	}
	assert.Equal(t, "test", Query[string](c, "a"))
	assert.Equal(t, 1, Query[int](c, "b"))
	assert.Equal(t, true, Query[bool](c, "c"))
	assert.Equal(t, float32(1.99), Query[float32](c, "d"))
	assert.Equal(t, 1.99, Query[float64](c, "d"))
}
