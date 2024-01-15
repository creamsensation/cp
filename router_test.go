package cp

import (
	"net/http"
	"slices"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/gox"
)

func TestRouter(t *testing.T) {
	cfg := config.Config{
		Languages: config.Languages{
			"cs": config.Language{Enabled: true, Default: true},
			"en": config.Language{Enabled: true, Default: false},
		},
		Router: config.Router{Localized: true},
		Security: config.Security{
			Csrf: config.SecurityCsrf{Enabled: true},
			Firewall: map[string]config.SecurityFirewall{
				naming.Main: {
					Enabled:     true,
					Controllers: []string{"test"},
				},
			},
			RateLimit: config.SecurityRateLimit{Enabled: true},
		},
	}
	cr := &core{config: cfg}
	r := createRouter(cr)
	cr.router = r
	cr.Routes(
		Route(
			Get(Map{"cs": "/domu", "en": "/home"}), Name("home"), Handler(
				func(c Control) Result {
					return gox.Text("home")
				},
			),
		),
	)
	t.Run(
		"prepare default language", func(t *testing.T) {
			assert.Equal(t, "cs", r.defaultLanguage)
		},
	)
	t.Run(
		"prepare localized path matcher", func(t *testing.T) {
			assert.True(t, slices.Contains([]string{`^/(cs|en)\b`, `^/(en|cs)\b`}, r.localizedPathMatcher.String()))
		},
	)
	t.Run(
		"prepare middlewares", func(t *testing.T) {
			assert.Equal(t, 2, len(r.middlewares))
		},
	)
	t.Run(
		"prepare firewalls", func(t *testing.T) {
			assert.Equal(t, 1, len(r.firewalls))
		},
	)
	t.Run(
		"prepare routes", func(t *testing.T) {
			r.prepareRoutes()
			langs := maps.Keys(r.localizedRoutes)
			assert.True(t, slices.Contains(langs, "cs"))
			assert.True(t, slices.Contains(langs, "en"))
			csLocalizedRoutes := r.localizedRoutes["cs"]
			enLocalizedRoutes := r.localizedRoutes["en"]
			assert.True(t, len(csLocalizedRoutes) == 1)
			assert.True(t, len(enLocalizedRoutes) == 1)
			assert.True(t, slices.Contains(csLocalizedRoutes[0].Method, http.MethodGet))
			assert.True(t, slices.Contains(enLocalizedRoutes[0].Method, http.MethodGet))
			assert.Equal(t, `^/cs/domu/?$`, csLocalizedRoutes[0].Matcher.String())
			assert.Equal(t, `^/en/home/?$`, enLocalizedRoutes[0].Matcher.String())
		},
	)
}
