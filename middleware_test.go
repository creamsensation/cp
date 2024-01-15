package cp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/constant/contentType"
	"github.com/creamsensation/cp/internal/constant/header"
	"github.com/creamsensation/cp/internal/tests"
	"github.com/creamsensation/form"
)

func TestMiddleware(t *testing.T) {
	cfg := config.Config{
		Cache: config.Cache{Adapter: cacheAdapter.Redis},
		Security: config.Security{
			Csrf: config.SecurityCsrf{
				Enabled:  true,
				Duration: time.Hour,
			},
		},
	}
	cr := &core{
		config: cfg,
		redis:  tests.CreateRedisConnection(t),
		ui:     createUi(),
	}
	cr.router = createRouter(cr)
	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)
	t.Run(
		"csrf", func(t *testing.T) {
			testCsrf := createCsrf(createControl(cr, req, httptest.NewRecorder()))
			handler := createCsrfMiddleware()
			token := testCsrf.Create("test", "localhost", "")
			ctrl := createControl(
				cr,
				httptest.NewRequest(
					http.MethodPost,
					"/test",
					strings.NewReader(
						fmt.Sprintf(
							"%s=%s&%s=%s",
							form.CsrfToken, token,
							form.CsrfName, "test",
						),
					),
				),
				httptest.NewRecorder(),
			)
			ctrl.request.Header.Add(header.ContentType, contentType.Form)
			assert.Nil(t, handler(ctrl))
		},
	)
	t.Run(
		"rate limit", func(t *testing.T) {
			handler := createRateLimitMiddleware(
				config.Security{
					RateLimit: config.SecurityRateLimit{
						Enabled:  true,
						Attempts: 5,
						Interval: time.Minute,
					},
				},
			)
			createTestControl := func() *control {
				return createControl(
					cr,
					httptest.NewRequest(
						http.MethodPost,
						"/test",
						nil,
					),
					httptest.NewRecorder(),
				)
			}
			assert.True(t, handler(createTestControl()) == nil)
			assert.True(t, handler(createTestControl()) == nil)
			assert.True(t, handler(createTestControl()) == nil)
			assert.True(t, handler(createTestControl()) == nil)
			assert.True(t, handler(createTestControl()) == nil)
			assert.True(t, handler(createTestControl()) != nil)
		},
	)
}
