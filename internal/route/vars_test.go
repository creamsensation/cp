package route

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestVars(t *testing.T) {
	t.Run(
		"create placeholders", func(t *testing.T) {
			assert.Equal(t, createVarsPlaceholders("/{lang:cs}"), map[string]string{"lang": "{lang:cs}"})
			assert.Equal(
				t, createVarsPlaceholders("/{lang:en}/users/{id:[0-9]+}"),
				map[string]string{"lang": "{lang:en}", "id": "{id:[0-9]+}"},
			)
		},
	)
}
