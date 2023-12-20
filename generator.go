package cp

import (
	"fmt"
	"strings"
	
	"github.com/creamsensation/cp/internal/constant/componentState"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/expiration"
	"github.com/creamsensation/cp/internal/constant/queryKey"
	"github.com/creamsensation/cp/internal/constant/requestVar"
	"github.com/creamsensation/cp/internal/querystring"
	"github.com/creamsensation/cp/internal/route"
)

type Generator interface {
	Link() LinkGenerator
	Query(arg Map) string
	Url() UrlGenerator
}

type LinkGenerator interface {
	Action(action string, arg ...Map) string
	Name(name string, arg ...Map) string
	SwitchLang(langCode string, overwrite ...Map) string
}

type UrlGenerator interface {
	Asset(path string) string
}

type generator struct {
	control *control
}

const (
	linkLevelDivider = "_"
)

func (g generator) Action(action string, arg ...Map) string {
	if g.control.component == nil {
		return action
	}
	action = g.control.component.Name() + linkLevelDivider + action
	shouldHaveModule := strings.Count(action, linkLevelDivider) == 3
	shouldHaveController := strings.Count(action, linkLevelDivider) >= 2
	shouldAddModulePrefix := len(g.control.route.Module) > 0 && !shouldHaveModule
	shouldAddControllerPrefix := len(g.control.route.Controller) > 0 && !shouldHaveController
	if shouldAddControllerPrefix {
		action = g.control.route.Controller + linkLevelDivider + action
	}
	if shouldAddModulePrefix {
		action = g.control.route.Module + linkLevelDivider + action
	}
	qm := Map{queryKey.Action: action}
	isCache := g.control.Config().Component.State == componentState.Cache
	args := make(Map)
	if len(arg) > 0 {
		args = arg[0]
		for k, v := range args {
			qm[k] = v
		}
	}
	if isCache {
		return g.control.Request().Path() + g.control.Generate().Query(qm)
	}
	queryStateParams := querystring.New(g.control.component).
		IgnoreInterface(componentControlInterfaceName).
		Override(args).
		Encode()
	if len(queryStateParams) == 0 {
		return g.control.Request().Path() + g.control.Generate().Query(qm)
	}
	return g.control.Request().Path() + g.control.Generate().Query(Map{queryKey.Action: action}) + "&" + queryStateParams
}

func (g generator) Asset(path string) string {
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}
	return fmt.Sprintf("%s/%s", g.control.config.Assets.PublicPath, path)
}

func (g generator) Link() LinkGenerator {
	return g
}

func (g generator) Name(name string, arg ...Map) string {
	if len(name) == 0 {
		return ""
	}
	if g.control.Request().Is().Localized() {
		if len(arg) == 0 {
			arg = make([]Map, 1)
			arg[0] = Map{requestVar.Lang: g.control.Request().Lang()}
		}
		if len(arg) > 0 {
			arg[0][requestVar.Lang] = g.control.Request().Lang()
		}
		localizedRoutes, ok := g.control.core.router.localizedRoutes[g.control.Request().Lang()]
		if !ok {
			return name
		}
		link, ok := g.generateLink(localizedRoutes, name, arg...)
		if ok {
			return link
		}
	}
	link, _ := g.generateLink(g.control.core.router.routes, name, arg...)
	return link
}

func (g generator) Query(arg Map) string {
	if len(arg) == 0 {
		return ""
	}
	result := make([]string, 0)
	for k, v := range arg {
		if v == nil {
			continue
		}
		result = append(result, fmt.Sprintf("%s=%v", k, v))
	}
	return "?" + strings.Join(result, "&")
}

func (g generator) SwitchLang(langCode string, overwrite ...Map) string {
	if !g.control.core.router.localized {
		return g.control.Request().Path()
	}
	var vars Map
	if !g.control.core.router.localized {
		return fmt.Sprintf("<%s:router is not localized>", langCode)
	}
	currentLang := Var[string](g.control, requestVar.Lang)
	index := -1
	for lc, lr := range g.control.core.router.localizedRoutes {
		if lc != currentLang {
			continue
		}
		for i, r := range lr {
			if r.Matcher.MatchString(g.control.Request().Path()) {
				index = i
				break
			}
		}
	}
	if index == -1 {
		return fmt.Sprintf("[%s:localized route does not exist>", langCode)
	}
	langRoutes, ok := g.control.core.router.localizedRoutes[langCode]
	if !ok {
		return fmt.Sprintf("<%s:language does not exist>", langCode)
	}
	if len(langRoutes)-1 < index {
		return fmt.Sprintf("/%s", langCode)
	}
	match := langRoutes[index]
	path := match.Path
	if len(overwrite) > 0 {
		vars = overwrite[0]
	}
	for vk, vp := range match.VarsPlaceholders {
		if vk == requestVar.Lang {
			path = strings.Replace(path, vp, langCode, 1)
			continue
		}
		v, ok := vars[vk]
		if ok {
			path = strings.Replace(path, vp, fmt.Sprintf("%v", v), 1)
		}
		if !ok {
			path = strings.Replace(path, vp, Var[string](g.control, vk), 1)
		}
	}
	g.control.Cookie().Set(cookieName.Lang, langCode, expiration.Lang)
	return path
}

func (g generator) Url() UrlGenerator {
	return g
}

func (g generator) generateLink(routes []route.Route, name string, args ...Map) (string, bool) {
	shouldHaveModule := strings.Count(name, linkLevelDivider) == 2
	shouldHaveController := strings.Count(name, linkLevelDivider) >= 1
	shouldAddModulePrefix := len(g.control.route.Module) > 0 && !shouldHaveModule
	shouldAddControllerPrefix := len(g.control.route.Controller) > 0 && !shouldHaveController
	for _, rt := range routes {
		if !strings.HasSuffix(rt.Name, name) {
			continue
		}
		if shouldAddControllerPrefix && rt.Controller == g.control.route.Controller {
			name = rt.Controller + linkLevelDivider + name
		}
		if shouldAddModulePrefix && rt.Module == g.control.route.Module {
			name = rt.Module + linkLevelDivider + name
		}
		if name != rt.Name {
			continue
		}
		if len(args) == 0 {
			return rt.Path, true
		}
		if len(rt.VarsPlaceholders) == 0 {
			return rt.Path, true
		}
		arg := args[0]
		for routeVarName, routeVarPlaceholder := range rt.VarsPlaceholders {
			if v, ok := arg[routeVarName]; ok {
				rt.Path = strings.Replace(rt.Path, routeVarPlaceholder, fmt.Sprintf("%v", v), 1)
			}
		}
		return rt.Path, true
	}
	return "/" + name, false
}
