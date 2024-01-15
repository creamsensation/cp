package cp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/translator"
)

func TestErrorHandler(t *testing.T) {
	t.Run(
		"check", func(t *testing.T) {
			e := createErrorHandler(&control{})
			err := errors.New("test error")
			defer func() {
				assert.Equal(t, err, recover())
			}()
			e.Check(err)
		},
	)
	t.Run(
		"throw", func(t *testing.T) {
			e := createErrorHandler(
				&control{
					core: &core{
						router:     createRouter(&core{}),
						translator: translator.New(t.TempDir()),
					},
					request: httptest.NewRequest(http.MethodGet, "/test", nil),
				},
			)
			err := errors.New("test error")
			defer func() {
				assert.Equal(t, err, recover())
			}()
			e.Message(err.Error()).Throw()
		},
	)
}
