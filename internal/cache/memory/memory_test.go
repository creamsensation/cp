package memory

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	m := &memory{
		data: make(map[string]memoryData),
		dir:  getDir(t.TempDir()),
	}
	
	t.Run(
		"set", func(t *testing.T) {
			assert.Nil(t, m.Set("test", "test", time.Minute))
		},
	)
	
	t.Run(
		"get", func(t *testing.T) {
			assert.Equal(t, "test", m.Get("test"))
		},
	)
	
	t.Run(
		"destroy", func(t *testing.T) {
			assert.Nil(t, m.Destroy("test"))
		},
	)
	
	t.Run(
		"restore", func(t *testing.T) {
			assert.Nil(t, m.Set("test", "test", time.Minute))
			m.data = make(map[string]memoryData)
			m.load()
			assert.Equal(t, true, m.Exists("test"))
		},
	)
}
