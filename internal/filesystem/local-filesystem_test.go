package filesystem

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestLocalFilesystem(t *testing.T) {
	fs := CreateLocalFilesystem(t.TempDir())
	t.Run(
		"create", func(t *testing.T) {
			assert.Nil(t, fs.Create("test/test.txt", []byte(`test string`)))
		},
	)
	t.Run(
		"read", func(t *testing.T) {
			fb, err := fs.Read("test/test.txt")
			assert.Nil(t, err)
			assert.Equal(t, "test string", string(fb))
		},
	)
	t.Run(
		"remove", func(t *testing.T) {
			assert.Nil(t, fs.Remove("test/test.txt"))
		},
	)
}
