package cp

import (
	"fmt"
	"regexp"
	"strings"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/requestVar"
	"github.com/creamsensation/cp/internal/firewall"
	"github.com/creamsensation/cp/internal/handler"
	"github.com/creamsensation/cp/internal/route"
)

type router struct {
	builders             []*route.Builder
	core                 *core
	config               config.Config
	defaultLanguage      string
	firewalls            map[string]*firewall.Firewall
	localized            bool
	localizedRoutes      map[string][]route.Route
	localizedPathMatcher *regexp.Regexp
	middlewares          []handler.Fn
	routes               []route.Route
}

func createRouter(core *core) *router {
	r := &router{
		builders:        make([]*route.Builder, 0),
		core:            core,
		config:          core.config,
		firewalls:       make(map[string]*firewall.Firewall),
		localized:       isRouterLocalized(core.config.Languages),
		localizedRoutes: make(map[string][]route.Route),
		middlewares:     make([]handler.Fn, 0),
	}
	r.onInit()
	return r
}

func (r *router) onInit() {
	r.prepareFirewalls()
	r.prepareMiddlewares()
	r.prepareDefaultLanguage()
	r.prepareLocalizedPathMatcher()
}

func (r *router) createServerHandler() *serverHandler {
	if r.localized {
		r.localizedRoutes = prepareLocalizedRoutes(r.builders, r.firewalls)
	}
	r.routes = prepareRoutes(r.builders, r.firewalls)
	return createServerHandler(r.core, r.routes)
}

func (r *router) prepareFirewalls() {
	if len(r.config.Security.Firewall) == 0 {
		return
	}
	for name, item := range r.config.Security.Firewall {
		matchers := make([]*regexp.Regexp, len(item.Patterns))
		for i, p := range item.Patterns {
			matchers[i] = regexp.MustCompile(p)
		}
		r.firewalls[name] = &firewall.Firewall{
			Enabled:       item.Enabled,
			Invert:        item.Invert,
			Modules:       item.Modules,
			Controllers:   item.Controllers,
			Routes:        item.Routes,
			Patterns:      item.Patterns,
			Matchers:      matchers,
			RedirectRoute: item.RedirectRoute,
			Roles:         item.Roles,
			Secret:        item.Secret,
		}
	}
}

func (r *router) prepareDefaultLanguage() {
	if !r.localized {
		return
	}
	for code, l := range r.config.Languages {
		if l.Enabled && l.Default {
			r.defaultLanguage = code
			break
		}
	}
}

func (r *router) prepareLocalizedPathMatcher() {
	if !r.localized {
		return
	}
	languages := make([]string, 0)
	for code, l := range r.config.Languages {
		if !l.Enabled {
			continue
		}
		languages = append(languages, code)
	}
	isPrefix := len(r.config.Router.PathPrefix) > 0
	if !isPrefix || (isPrefix && !r.config.Router.PreferPrefix) {
		r.localizedPathMatcher = regexp.MustCompile(fmt.Sprintf(`^/(%s)\b`, strings.Join(languages, "|")))
		return
	}
	if isPrefix && r.config.Router.PreferPrefix {
		r.localizedPathMatcher = regexp.MustCompile(
			fmt.Sprintf(
				`^/%s/(%s)\b`, r.config.Router.PathPrefix, strings.Join(languages, "|"),
			),
		)
	}
}

func (r *router) prepareMiddlewares() {
	if r.config.Security.Csrf.Enabled {
		r.middlewares = append(
			r.middlewares,
			createCsrfMiddleware(),
		)
	}
	if r.config.Security.RateLimit.Enabled {
		r.middlewares = append(
			r.middlewares,
			createRateLimitMiddleware(r.config.Security),
		)
	}
}

func prepareRoutes(items []*route.Builder, firewalls map[string]*firewall.Firewall) []route.Route {
	result := make([]route.Route, 0)
	for i, item := range items {
		if len(item.Controller) > 0 {
			items[i].Name = item.Controller + linkLevelDivider + items[i].Name
		}
		if len(item.Module) > 0 {
			items[i].Name = item.Module + linkLevelDivider + items[i].Name
		}
	}
	for i, item := range items {
		if !item.Ok {
			continue
		}
		item.Route.Firewalls = make(map[string]*firewall.Route)
		for firewallName, f := range firewalls {
			var exist bool
			for _, m := range f.Modules {
				if item.Module == m {
					items[i].Firewalls[firewallName] = createRouteFirewall(f)
					exist = true
				}
			}
			if exist {
				continue
			}
			for _, r := range f.Controllers {
				if item.Controller == r {
					items[i].Firewalls[firewallName] = createRouteFirewall(f)
					exist = true
				}
			}
			if exist {
				continue
			}
			for _, name := range f.Routes {
				if item.Name == name {
					items[i].Firewalls[firewallName] = createRouteFirewall(f)
					exist = true
				}
			}
			if exist {
				continue
			}
			for _, m := range f.Matchers {
				if m.MatchString(item.Route.Path) {
					items[i].Firewalls[firewallName] = createRouteFirewall(f)
				}
			}
		}
		result = append(result, item.Route)
	}
	return result
}

func prepareLocalizedRoutes(items []*route.Builder, firewalls map[string]*firewall.Firewall) map[string][]route.Route {
	result := make(map[string][]route.Route)
	for langCode, b := range items {
		for i, item := range b.LocalizedRoute {
			if len(item.Controller) > 0 {
				items[langCode].LocalizedRoute[i].Name = item.Controller + linkLevelDivider + items[langCode].LocalizedRoute[i].Name
			}
			if len(item.Module) > 0 {
				items[langCode].LocalizedRoute[i].Name = item.Module + linkLevelDivider + items[langCode].LocalizedRoute[i].Name
			}
		}
	}
	for _, item := range items {
		for langCode, localizedRoute := range item.LocalizedRoute {
			if !localizedRoute.Ok {
				continue
			}
			localizedRoute.Firewalls = make(map[string]*firewall.Route)
			for firewallName, f := range firewalls {
				var exist bool
				for _, m := range f.Modules {
					if item.LocalizedRoute[langCode].Module == m {
						item.LocalizedRoute[langCode].Firewalls[firewallName] = createRouteFirewall(f)
						exist = true
					}
				}
				if exist {
					continue
				}
				for _, r := range f.Controllers {
					if item.LocalizedRoute[langCode].Controller == r {
						item.LocalizedRoute[langCode].Firewalls[firewallName] = createRouteFirewall(f)
						exist = true
					}
				}
				if exist {
					continue
				}
				for _, name := range f.Routes {
					if item.LocalizedRoute[langCode].Name == name {
						item.LocalizedRoute[langCode].Firewalls[firewallName] = createRouteFirewall(f)
						exist = true
					}
				}
				if exist {
					continue
				}
				for _, m := range f.Matchers {
					if m.MatchString(
						strings.Replace(
							item.LocalizedRoute[langCode].Route.Path,
							fmt.Sprintf(`{%s:%s}`, requestVar.Lang, langCode),
							langCode,
							1,
						),
					) {
						item.LocalizedRoute[langCode].Firewalls[firewallName] = createRouteFirewall(f)
					}
				}
			}
			if _, ok := result[langCode]; !ok {
				result[langCode] = make([]route.Route, 0)
			}
			result[langCode] = append(result[langCode], localizedRoute.Route)
		}
	}
	return result
}

func isRouterLocalized(languages config.Languages) bool {
	if len(languages) == 0 {
		return false
	}
	for _, l := range languages {
		if l.Enabled {
			return true
		}
	}
	return false
}

func createRouteFirewall(f *firewall.Firewall) *firewall.Route {
	return &firewall.Route{
		Enabled:       f.Enabled,
		Invert:        f.Invert,
		RedirectRoute: f.RedirectRoute,
		Roles:         f.Roles,
		Secret:        f.Secret,
	}
}
