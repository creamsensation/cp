package route

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	t.Run(
		"create matcher", func(t *testing.T) {
			assert.Equal(t, `^/test/?$`, createMatcher("/test").String())
			assert.Equal(t, `^/test/abc/?$`, createMatcher("/test/abc").String())
			assert.Equal(t, `^/test/abc/?$`, createMatcher("test/abc/").String())
			assert.Equal(t, `^/test/[0-9]+/?$`, createMatcher("test/{id:[0-9]+}").String())
			assert.Equal(t, `^/.*/?$`, createMatcher("*").String())
		},
	)
}
