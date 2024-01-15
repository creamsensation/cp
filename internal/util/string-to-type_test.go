package util

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestStringToType(t *testing.T) {
	assert.Equal(t, "test", StringToType[string]("test"))
	assert.Equal(t, 1, StringToType[int]("1"))
	assert.Equal(t, true, StringToType[bool]("true"))
	assert.Equal(t, float32(1.99), StringToType[float32]("1.99"))
	assert.Equal(t, 1.99, StringToType[float64]("1.99"))
}
