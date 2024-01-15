package route

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestUtil(t *testing.T) {
	t.Run(
		"flat", func(t *testing.T) {
			b := CreateBuilder(
				CreatePathConfig("/test"),
			).Group(
				CreateBuilder(
					CreatePathConfig("/1"),
				).Group(
					CreateBuilder(
						CreatePathConfig("/2"),
					),
				),
			)
			flatten := CreateFlatBuilders(b)
			assert.Equal(t, 3, len(flatten))
		},
	)
	t.Run(
		"is localized", func(t *testing.T) {
			b := CreateBuilder(
				CreatePathConfig(map[string]string{"cs": "/kosik", "en": "/cart"}),
			)
			assert.True(t, true, IsLocalized(b))
		},
	)
}
