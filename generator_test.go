package cp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/tests"
)

func TestGenerator(t *testing.T) {
	testCore := &core{
		config: config.Config{
			Languages: tests.Languages,
			Router:    config.Router{Localized: true},
		},
	}
	testCore.router = createRouter(testCore)
	testCore.Routes(
		Route(
			Get(""),
			Name("home"),
		),
		Route(
			Get("blog"),
			Name("blog"),
		).Group(
			Route(
				Get("{id:[0-9]+}"),
				Name("detail"),
			),
		),
	)
	testCore.router.prepareRoutes()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := createControl(testCore, req, res)
	t.Run(
		"generate link", func(t *testing.T) {
			assert.Equal(t, "/", c.Link("home"))
			assert.Equal(t, "/blog", c.Link("blog"))
			assert.Equal(t, "/blog/1", c.Link("blog-detail", Map{"id": 1}))
		},
	)
	t.Run(
		"generate link with query", func(t *testing.T) {
			q := c.Link("blog") + c.Generate().Query(Map{"id": 1, "name": "Dominik"})
			assert.Contains(t, q, "/blog")
			assert.Contains(t, q, "id=1")
			assert.Contains(t, q, "name=Dominik")
		},
	)
}

func TestLocalizedGenerator(t *testing.T) {
	testCore := &core{
		config: config.Config{
			Languages: tests.Languages,
			Router:    config.Router{Localized: true},
		},
	}
	testCore.router = createRouter(testCore)
	testCore.Routes(
		Route(
			Get(Map{"cs": "clanky", "en": "articles"}),
			Name("articles"),
		),
	)
	testCore.router.prepareRoutes()
	req := httptest.NewRequest(http.MethodGet, "/cs/clanky", nil)
	res := httptest.NewRecorder()
	c := createControl(testCore, req, res)
	createLifecycle(c).run()
	t.Run(
		"generate link", func(t *testing.T) {
			assert.Equal(t, "/cs/clanky", c.Link("articles"))
			assert.Equal(t, "/cs/clanky", c.Generate().Link().Name("articles"))
			assert.Equal(t, "/en/articles", c.Generate().Link().SwitchLang("en"))
		},
	)
}
