package assets

import (
	"fmt"
	"slices"
	"strings"
	
	"github.com/creamsensation/gox"
)

type Assets interface {
	Get(path string) string
	GetStyles() gox.Node
	GetScripts() gox.Node
	AddStyle(path string) Assets
	AddScript(path string, afterDomLoaded bool) Assets
}

type assets struct {
	publicPath   string
	styles       []string
	scripts      []string
	scriptsDefer []string
}

func New(publicPath string, styles, scripts, scriptsDefer []string) Assets {
	if !strings.HasPrefix(publicPath, "/") {
		publicPath = "/" + publicPath
	}
	return &assets{
		publicPath:   publicPath,
		styles:       styles,
		scripts:      scripts,
		scriptsDefer: scriptsDefer,
	}
}

func (a *assets) Get(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return fmt.Sprintf("%s/%s", a.publicPath, path)
}

func (a *assets) GetStyles() gox.Node {
	return gox.Range(
		a.styles, func(value string, _ int) gox.Node {
			if !strings.HasPrefix(value, "/") {
				value = "/" + value
			}
			return gox.Link(
				gox.Rel("stylesheet"),
				gox.Type("text/css"),
				gox.Href(value),
			)
		},
	)
}

func (a *assets) GetScripts() gox.Node {
	return gox.Fragment(
		gox.Range(
			a.scripts, func(value string, _ int) gox.Node {
				if !strings.HasPrefix(value, "/") {
					value = "/" + value
				}
				return gox.Script(gox.Src(value), gox.Type("module"))
			},
		),
		gox.Range(
			a.scriptsDefer, func(value string, _ int) gox.Node {
				if !strings.HasPrefix(value, "/") {
					value = "/" + value
				}
				return gox.Script(gox.Defer(), gox.Src(value), gox.Type("module"))
			},
		),
	)
}

func (a *assets) AddStyle(path string) Assets {
	if slices.Contains(a.styles, path) {
		return a
	}
	a.styles = append(a.styles, path)
	return a
}

func (a *assets) AddScript(path string, afterDomLoaded bool) Assets {
	if !afterDomLoaded && slices.Contains(a.scripts, path) {
		return a
	}
	if afterDomLoaded && slices.Contains(a.scriptsDefer, path) {
		return a
	}
	if !afterDomLoaded {
		a.scripts = append(a.scripts, path)
	}
	if afterDomLoaded {
		a.scriptsDefer = append(a.scriptsDefer, path)
	}
	return a
}
