package util

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestGetInterfaceName(t *testing.T) {
	type testInterface interface{}
	assert.Equal(t, "util.testInterface", GetInterfaceName[testInterface]())
}
