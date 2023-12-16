package route

import (
	"fmt"
	"slices"
	"strings"
	
	cfg "github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/requestVar"
	"github.com/creamsensation/cp/internal/handler"
)

func Process(b *Builder, parent *Builder, languages cfg.Languages, rc cfg.Router) {
	if !rc.PreferPrefix {
		createBuilderPathPrefix(b, rc)
	}
	if parent != nil && parent.Localized && !b.Localized {
		createNestedBuilderLocalized(b, parent)
	}
	if b.Localized {
		createLocalizedRouteBuilders(b)
		processLocalizedRoute(b, parent, rc)
	}
	if !b.Localized {
		processRoute(b, parent, languages, rc)
	}
	if len(b.Subroutes) > 0 {
		for _, sub := range b.Subroutes {
			Process(sub, b, languages, rc)
		}
	}
}

func processLocalizedRoute(b *Builder, parent *Builder, rc cfg.Router) {
	for _, item := range b.Configs {
		switch c := item.(type) {
		case config[string]:
			switch c.configType {
			case ConfigName:
				for _, lb := range b.LocalizedRoute {
					lb.Name = c.value
				}
			case ConfigLayout:
				for _, lb := range b.LocalizedRoute {
					lb.Layout = c.value
				}
			}
		case config[handler.Fn]:
			for _, lb := range b.LocalizedRoute {
				lb.Fn = c.value
			}
		default:
			continue
		}
	}
	for _, item := range b.Configs {
		switch c := item.(type) {
		case config[pathConfig]:
			switch paths := c.value.path.(type) {
			case map[string]any:
				for langCode, localizedPath := range paths {
					path := fmt.Sprintf("%v", localizedPath)
					if len(c.value.method) > 0 {
						for _, mtd := range c.value.method {
							if !slices.Contains(b.LocalizedRoute[langCode].Method, mtd) {
								b.LocalizedRoute[langCode].Method = append(b.LocalizedRoute[langCode].Method, mtd)
							}
						}
					}
					// if len(path) > 0 {
					if parent == nil {
						path = prefixPathWithLangVar(path, langCode)
					}
					preparePath(b.LocalizedRoute[langCode], path)
					if parent != nil {
						prepareSubroutePrefixes(b.LocalizedRoute[langCode], parent.LocalizedRoute[langCode])
					}
					if rc.PreferPrefix {
						preparePathPrefix(b.LocalizedRoute[langCode], rc)
					}
					prepareRoute(b.LocalizedRoute[langCode])
				}
				// }
			case nil:
				for _, langCode := range getLocalizedRouteLanguages(b) {
					for _, mtd := range c.value.method {
						if !slices.Contains(b.LocalizedRoute[langCode].Method, mtd) {
							b.LocalizedRoute[langCode].Method = append(b.LocalizedRoute[langCode].Method, mtd)
						}
					}
				}
			}
		default:
			continue
		}
	}
	for _, langCode := range getLocalizedRouteLanguages(b) {
		b.LocalizedRoute[langCode].Module = b.Module
		b.LocalizedRoute[langCode].Controller = b.Controller
	}
}

func processRoute(b *Builder, parent *Builder, languages cfg.Languages, rc cfg.Router) {
	for _, item := range b.Configs {
		switch c := item.(type) {
		case config[string]:
			switch c.configType {
			case ConfigName:
				b.Name = c.value
			case ConfigLayout:
				b.Layout = c.value
			}
		case config[bool]:
			switch c.configType {
			case ConfigLocalize:
				createBuilderLocalizedWithLanguages(b, languages)
			}
		case config[handler.Fn]:
			b.Fn = c.value
		default:
			continue
		}
	}
	if b.Localized {
		createLocalizedRouteBuilders(b)
		processLocalizedRoute(b, parent, rc)
		return
	}
	
	for _, item := range b.Configs {
		switch c := item.(type) {
		case config[pathConfig]:
			switch path := c.value.path.(type) {
			case string:
				// if len(path) > 0 {
				preparePath(b, path)
				if parent != nil {
					prepareSubroutePrefixes(b, parent)
				}
				if rc.PreferPrefix {
					preparePathPrefix(b, rc)
				}
				prepareRoute(b)
				// }
				if len(c.value.method) > 0 {
					for _, mtd := range c.value.method {
						if !slices.Contains(b.Method, mtd) {
							b.Method = append(b.Method, mtd)
						}
					}
				}
			case nil:
				if len(c.value.method) > 0 {
					for _, mtd := range c.value.method {
						if !slices.Contains(b.Method, mtd) {
							b.Method = append(b.Method, mtd)
						}
					}
				}
			}
		default:
			continue
		}
	}
}

func processMethods(methods ...Config) []string {
	result := make([]string, len(methods))
	for i, item := range methods {
		switch c := item.(type) {
		case config[pathConfig]:
			result[i] = c.value.method[0]
		}
	}
	return result
}

func prepareSubroutePrefixes(b *Builder, parent *Builder) {
	b.Path = parent.Path + b.Path
	b.Name = parent.Name + "-" + b.Name
}

func preparePath(b *Builder, path string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") && path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	b.Path = path
}

func preparePathPrefix(b *Builder, rc cfg.Router) {
	if len(rc.PathPrefix) == 0 {
		return
	}
	b.Path = strings.TrimPrefix(b.Path, "/")
	b.Path = "/" + rc.PathPrefix + "/" + b.Path
}

func createBuilderPathPrefix(b *Builder, rc cfg.Router) {
	if len(rc.PathPrefix) == 0 {
		return
	}
	for i, bc := range b.Configs {
		switch c := bc.(type) {
		case config[pathConfig]:
			switch p := c.value.path.(type) {
			case map[string]any:
				for lc := range p {
					path := fmt.Sprintf("%v", p[lc])
					path = strings.TrimPrefix(path, "/")
					p[lc] = "/" + rc.PathPrefix + "/" + path
				}
				c.value.path = p
				b.Configs[i] = c
			case string:
				p = strings.TrimPrefix(p, "/")
				p = "/" + rc.PathPrefix + "/" + p
				c.value.path = p
				b.Configs[i] = c
			}
		}
	}
}

func prepareRoute(r *Builder) {
	r.VarsPlaceholders = createVarsPlaceholders(r.Path)
	r.Matcher = createMatcher(r.Path)
	r.Ok = true
}

func prefixPathWithLangVar(path string, langCode string) string {
	if len(langCode) == 0 {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = "/" + fmt.Sprintf("{%s:%s}", requestVar.Lang, langCode) + path
	return path
}

func getLocalizedRouteLanguages(b *Builder) []string {
	result := make([]string, 0)
	for _, item := range b.Configs {
		switch c := item.(type) {
		case config[pathConfig]:
			switch paths := c.value.path.(type) {
			case map[string]any:
				for langCode := range paths {
					result = append(result, langCode)
				}
			}
		}
	}
	return result
}

func createLocalizedRouteBuilders(b *Builder) {
	for _, item := range b.Configs {
		switch c := item.(type) {
		case config[pathConfig]:
			switch paths := c.value.path.(type) {
			case map[string]any:
				for langCode := range paths {
					b.LocalizedRoute[langCode] = CreateBuilder(b.Configs...)
				}
			}
		}
	}
}

func createBuilderLocalizedWithLanguages(b *Builder, languages cfg.Languages) {
	for i, item := range b.Configs {
		cc := item
		switch c := item.(type) {
		case config[pathConfig]:
			localizedRoutes := make(map[string]any)
			switch path := c.value.path.(type) {
			case string:
				for langCode := range languages {
					localizedRoutes[langCode] = path
				}
			}
			c.value.path = localizedRoutes
			cc = c
		}
		b.Configs[i] = cc
	}
	b.Localized = true
}

func createNestedBuilderLocalized(b *Builder, parent *Builder) {
	for i, item := range b.Configs {
		cc := item
		switch c := item.(type) {
		case config[pathConfig]:
			localizedRoutes := make(map[string]any)
			switch path := c.value.path.(type) {
			case string:
				for langCode := range parent.LocalizedRoute {
					localizedRoutes[langCode] = path
				}
			}
			c.value.path = localizedRoutes
			cc = c
		}
		b.Configs[i] = cc
	}
	b.Localized = true
}
