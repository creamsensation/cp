package cp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/contentType"
	"github.com/creamsensation/cp/internal/constant/header"
)

func TestRequest(t *testing.T) {
	c := &core{
		config: config.Config{
			Languages: config.Languages{
				"cs": config.Language{Enabled: true, Default: true},
				"en": config.Language{Enabled: true, Default: false},
			},
			Router: config.Router{Localized: true},
		},
	}
	c.router = createRouter(c)
	c.router.onInit()
	c.router.prepareRoutes()
	t.Run(
		"content type", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, contentType.Form, ctrl.Request().ContentType())
		},
	)
	t.Run(
		"form", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("test=abc"))
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, "abc", ctrl.Request().Form().Value("test"))
		},
	)
	t.Run(
		"lang", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cs/test", nil)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			ctrl.vars["lang"] = "cs"
			assert.Equal(t, "cs", ctrl.Request().Lang())
		},
	)
	t.Run(
		"host & ip", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, "http://example.com", ctrl.Request().Host())
			assert.Equal(t, "localhost", ctrl.Request().Ip())
		},
	)
	t.Run(
		"method & path", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, http.MethodPost, ctrl.Request().Method())
			assert.Equal(t, "/test", ctrl.Request().Path())
		},
	)
	t.Run(
		"protocol", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, "http", ctrl.Request().Protocol())
		},
	)
	t.Run(
		"query", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test?test=abc", nil)
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, "abc", ctrl.Request().Query("test"))
		},
	)
	t.Run(
		"user-agent", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(header.UserAgent, "test")
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			assert.Equal(t, "test", ctrl.Request().UserAgent())
		},
	)
	t.Run(
		"var", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(header.ContentType, contentType.Form)
			ctrl := createControl(
				c,
				req,
				httptest.NewRecorder(),
			)
			ctrl.vars["test"] = "abc"
			assert.Equal(t, "abc", ctrl.Request().Var("test"))
		},
	)
}
