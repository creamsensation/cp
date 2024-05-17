package cp

import (
	"net/http"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	t.Run(
		"basic", func(t *testing.T) {
			result := "test"
			w := TestRoute(
				TestRouteParam{
					Method:  http.MethodGet,
					Path:    "/test/",
					TempDir: t.TempDir(),
					Handler: func(c Ctx) error {
						return c.Response().Text(result)
					},
				},
			)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, result, string(w.Body.Bytes()))
		},
	)
}
