package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	c := &control{
		vars: map[string]string{
			"a": "1",
			"b": "test",
			"c": "true",
			"d": "1.99",
			"e": "99.05",
		},
	}
	assert.Equal(t, 1, Var[int](c, "a"))
	assert.Equal(t, "test", Var[string](c, "b"))
	assert.Equal(t, true, Var[bool](c, "c"))
	assert.Equal(t, float32(1.99), Var[float32](c, "d"))
	assert.Equal(t, 99.05, Var[float64](c, "e"))
	
}
