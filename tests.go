package cp

import (
	"io"
	"net/http"
	"net/http/httptest"
	
	"github.com/creamsensation/cache/memory"
	"github.com/creamsensation/config"
)

type TestCtxParam struct {
	Config         config.Config
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	MatchedRoute   *Route
	Routes         *[]*Route
}

type TestRouteParam struct {
	Config  *config.Config
	Method  string
	Path    string
	TempDir string
	Body    io.Reader
	Handler func(c Ctx) error
}

func TestCtx(param TestCtxParam) Ctx {
	return createContext(
		ctxParam{
			config:       param.Config,
			matchedRoute: param.MatchedRoute,
			routes:       param.Routes,
			r:            param.Request,
			w:            param.ResponseWriter,
			layout:       new(layout),
		},
	)
}

func TestRoute(param TestRouteParam) *httptest.ResponseRecorder {
	cfg := config.Config{
		App: config.App{
			Name: "test",
		},
		Cache: config.Cache{
			Memory: memory.New(param.TempDir),
		},
		Router: config.Router{
			Recover: true,
		},
	}
	if param.Config != nil {
		cfg = *param.Config
	}
	c := New(cfg)
	c.Route(
		param.Path,
		param.Handler,
		Method(param.Method),
	)
	r := httptest.NewRequest(param.Method, param.Path, param.Body)
	w := httptest.NewRecorder()
	c.Mux().ServeHTTP(w, r)
	return w
}
