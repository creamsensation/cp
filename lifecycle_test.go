package cp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/gox"
)

func TestLifecycle(t *testing.T) {
	cfg := config.Config{
		Cache:     config.Cache{Adapter: cacheAdapter.Redis},
		Languages: config.Languages{"cs": {Enabled: true, Default: true}, "en": {Enabled: true}},
		Security: config.Security{
			Firewall: map[string]config.SecurityFirewall{
				naming.Main: {
					Enabled:  true,
					Patterns: []string{`^/protected\b`},
					Roles:    []string{"owner"},
				},
			},
			Role: map[string]config.SecurityRole{
				"owner": {Super: true},
			},
		},
	}
	cr := &core{
		config: cfg,
		ui:     createUi(),
	}
	cr.router = createRouter(cr)
	cr.Routes(
		Route(
			Get("/test"), Name("test"), Handler(
				func(c Control) Result {
					return c.Response().Render(gox.Text("test"))
				},
			),
		),
		Route(
			Get("/protected"), Name("protected"), Handler(
				func(c Control) Result {
					return c.Response().Render(gox.Text("protected"))
				},
			),
		),
	)
	cr.router.onInit()
	cr.router.prepareRoutes()
	t.Run(
		"non-localized", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			l := createLifecycle(ctrl)
			l.run()
			assert.Equal(t, "test", any(ctrl.response).(*httptest.ResponseRecorder).Body.String())
		},
	)
}
