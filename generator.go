package cp

import (
	"fmt"
	"net/url"
	"strings"
	
	"github.com/creamsensation/gox"
	
	"github.com/creamsensation/csrf"
	"github.com/creamsensation/form"
)

type Generator interface {
	Assets() gox.Node
	Action(name string, args ...Map) string
	Csrf(name string) gox.Node
	Link(name string, args ...Map) string
	PublicUrl(path string) string
	Query(args Map) string
	SwitchLang(langCode string) string
}

type generator struct {
	*ctx
}

func (g *generator) Assets() gox.Node {
	if g.assets == nil {
		return gox.Fragment()
	}
	return gox.Fragment(
		gox.Range(
			g.assets.styles, func(style string, _ int) gox.Node {
				return gox.Link(gox.Rel("stylesheet"), gox.Type("text/css"), gox.Href(style))
			},
		),
		gox.Range(
			g.assets.scripts, func(style string, _ int) gox.Node {
				return gox.Script(gox.Defer(), gox.Src(style))
			},
		),
	)
}

func (g *generator) Action(name string, args ...Map) string {
	if g.component == nil {
		return ""
	}
	qpm := Map{Action: g.route.Name + namePrefixDivider + g.component.name + namePrefixDivider + name}
	for k, vals := range g.Request().Raw().URL.Query() {
		for _, v := range vals {
			qpm[k] = v
		}
	}
	if len(args) > 0 {
		for k, v := range args[0] {
			qpm[k] = v
		}
	}
	return g.Request().Path() + g.Query(qpm)
}

func (g *generator) Csrf(name string) gox.Node {
	token := g.csrf.MustCreate(
		csrf.Token{
			Name:      name,
			Ip:        g.Request().Ip(),
			UserAgent: g.Request().UserAgent(),
		},
	)
	return form.Csrf(name, token)
}

func (g *generator) Link(name string, args ...Map) string {
	if !strings.Contains(name, namePrefixDivider) {
		name = g.route.Name + namePrefixDivider + name
	}
	l := g.Lang().Current()
	for _, r := range *g.routes {
		if g.config.Localization.Enabled && !g.config.Localization.Path {
			if r.Name == name {
				return g.replacePathParamsWithArgs(r.Path, args...)
			}
			continue
		}
		if g.config.Localization.Enabled && r.Name == name && r.Lang == l {
			return g.replacePathParamsWithArgs(r.Path, args...)
		}
		if r.Name == name {
			return g.replacePathParamsWithArgs(r.Path, args...)
		}
	}
	return ""
}

func (g *generator) Query(args Map) string {
	if len(args) == 0 {
		return ""
	}
	result := make([]string, 0)
	for k, v := range args {
		if v == nil {
			continue
		}
		result = append(result, fmt.Sprintf("%s=%v", k, v))
	}
	return "?" + strings.Join(result, "&")
}

func (g *generator) PublicUrl(path string) string {
	r, err := url.JoinPath(g.config.App.Public, path)
	if err != nil {
		return path
	}
	return r
}

func (g *generator) SwitchLang(langCode string) string {
	path := g.Request().Path()
	name := g.Request().Name()
	g.cookie.Set(langCookieKey, langCode, langCookieDuration)
	if !g.config.Localization.Path {
		return path
	}
	for _, r := range *g.routes {
		if r.Name == name && r.Lang == langCode {
			return r.Path
		}
	}
	return path
}

func (g *generator) replacePathParamsWithArgs(path string, args ...Map) string {
	if len(args) == 0 {
		return path
	}
	replace := make([]string, 0)
	for k, v := range args[0] {
		replace = append(replace, "{"+k+"}", fmt.Sprintf("%v", v))
	}
	r := strings.NewReplacer(replace...)
	return r.Replace(path)
}
