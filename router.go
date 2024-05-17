package cp

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	
	"github.com/creamsensation/config"
	"github.com/creamsensation/firewall"
)

type Router interface {
	Static(path, dir string) Router
	Route(path any, handler Handler, config ...RouteConfig) Router
	Group(path any, name ...string) Router
}

type router struct {
	core        *core
	config      config.Config
	mux         *http.ServeMux
	prefix      config.Prefix
	middlewares []Handler
	assets      *assets
	routes      *[]*Route
}

const (
	paramRegex = `[0-9a-zA-Z]+`
)

func (r *router) Static(path, dir string) Router {
	r.mux.Handle(http.MethodGet+" "+path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
	return r
}

func (r *router) Route(path any, fn Handler, config ...RouteConfig) Router {
	switch v := path.(type) {
	case string:
		if r.config.Localization.Path {
			for _, item := range r.config.Localization.Languages {
				r.createRoute(v, fn, item.Code, config...)
			}
		}
		if !r.config.Localization.Path {
			r.createRoute(v, fn, "", config...)
		}
	case map[string]string:
		for l, p := range v {
			r.createRoute(p, fn, l, config...)
		}
	}
	return r
}

func (r *router) Group(path any, name ...string) Router {
	var routerName string
	if len(name) > 0 {
		routerName = name[0]
	}
	return &router{
		core:   r.core,
		config: r.config,
		mux:    r.mux,
		prefix: config.Prefix{
			Path: r.mergePrefixPath(r.prefix.Path, path),
			Name: r.prefix.Name + routerName,
		},
		middlewares: r.middlewares,
		assets:      r.assets,
		routes:      r.routes,
	}
}

func (r *router) createGetWildcardRoute() {
	method := http.MethodGet
	path := "/{path...}"
	r.mux.HandleFunc(
		r.createRoutePattern(method, path),
		r.createHandler(
			method, path, "wildcard", func(c Ctx) error {
				c.Response().Status(http.StatusNotFound)
				return r.core.errorHandler(c)
			},
		),
	)
}

func (r *router) createRoute(path string, fn Handler, lang string, config ...RouteConfig) {
	var name string
	methods := make([]string, 0)
	for _, cfg := range config {
		switch cfg.Type {
		case routeMethod:
			methods = cfg.Value.([]string)
		case routeName:
			name = cfg.Value.(string)
		}
	}
	if len(methods) == 0 {
		methods = append(methods, httpMethods...)
	}
	if r.prefix.Path != nil {
		switch v := r.prefix.Path.(type) {
		case string:
			path = r.mustJoinPath(v, path)
		case map[string]string:
			lp, ok := v[lang]
			if ok {
				path = r.mustJoinPath(lp, path)
			}
		}
	}
	path = r.prefixPathWithLangIfEnabled(path, lang)
	if len(r.prefix.Name) > 0 {
		name = r.prefix.Name + namePrefixDivider + name
	}
	*r.routes = append(
		*r.routes, &Route{
			Lang:      lang,
			Path:      path,
			Name:      name,
			Methods:   methods,
			Matcher:   r.createMatcher(path),
			Firewalls: r.createFirewalls(path, name),
		},
	)
	for _, method := range methods {
		r.mux.HandleFunc(
			r.createRoutePattern(method, path),
			r.createHandler(method, path, name, fn),
		)
	}
}

func (r *router) createFirewalls(path, name string) []firewall.Firewall {
	result := make([]firewall.Firewall, 0)
	for _, f := range r.config.Security.Firewalls {
		if f.Match(path) {
			result = append(result, f)
			continue
		}
		if f.MatchPath(path) {
			result = append(result, f)
			continue
		}
		if f.MatchGroup(name) {
			result = append(result, f)
			continue
		}
	}
	return result
}

func (r *router) createRoutePattern(method, path string) string {
	return method + " " + r.formatPatternPath(path)
}

func (r *router) createMatcher(path string) *regexp.Regexp {
	parts := strings.Split(path, "/")
	res := make([]string, len(parts))
	for i, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			res[i] = paramRegex
			continue
		}
		res[i] = part
	}
	return regexp.MustCompile(strings.Join(res, "/"))
}

func (r *router) formatPatternPath(path string) string {
	if strings.Contains(path, "...") {
		return path
	}
	if !strings.HasSuffix(path, "/") {
		return path + "/{$}"
	}
	return path + "{$}"
}

func (r *router) mustJoinPath(basePath string, path string) string {
	if strings.Contains(path, "...") {
		return strings.TrimSuffix(basePath, "/") + "/" + strings.TrimPrefix(path, "/")
	}
	p, err := url.JoinPath(basePath, path)
	if err != nil {
		panic(err)
	}
	return p
}

func (r *router) mergePrefixPath(prefixPath any, path any) any {
	switch pp := prefixPath.(type) {
	case string:
		switch p := path.(type) {
		case string:
			return pp + p
		}
	case map[string]string:
		switch p := path.(type) {
		case map[string]string:
			for l, item := range pp {
				p[l] = item + p[l]
			}
			return p
		}
	}
	return path
}

func (r *router) prefixPathWithLangIfEnabled(path, lang string) string {
	if r.config.Localization.Path && !strings.HasSuffix(path, "/"+lang+"/") {
		return r.mustJoinPath("/"+lang+"/", path)
	}
	return path
}

func (r *router) createHandler(method, path, name string, fn Handler) func(http.ResponseWriter, *http.Request) {
	return handler{
		core:   r.core,
		method: method,
		path:   path,
		name:   name,
	}.create(fn)
}
