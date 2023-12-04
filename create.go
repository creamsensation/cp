package cp

import (
	"fmt"
	"net/http"
	"strings"
	
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/cp/internal/constant/componentState"
	"github.com/creamsensation/cp/internal/constant/queryKey"
	"github.com/creamsensation/cp/internal/constant/requestVar"
	"github.com/creamsensation/cp/internal/querystring"
	"github.com/creamsensation/cp/internal/route"
	"github.com/creamsensation/form"
	"github.com/creamsensation/gox"
	"github.com/creamsensation/hx"
)

type Create interface {
	Component(component component) gox.Node
	Defer(link string, nodes ...gox.Node) gox.Node
	Email() Email
	FormBuilder(fields ...*form.FieldBuilder) *form.Builder
	Link(name string, args ...map[string]any) string
	Pdf() Pdf
	Query(qm Map) string
}

type create struct {
	*control
	component component
}

const (
	linkLevelDivider = "_"
)

func (c create) Component(ct component) gox.Node {
	cc := createComponentControl(c.control, ct)
	cl := createComponentLifecycle(cc)
	cl.run()
	return cl.node()
}

func (c create) Defer(link string, nodes ...gox.Node) gox.Node {
	return gox.Div(
		hx.Get(link),
		hx.Trigger("load"),
		hx.Swap(hx.SwapOuterHtml),
		hx.Headers(map[string]any{hx.RequestHeaderTrigger: "load"}),
		gox.Fragment(nodes...),
	)
}

func (c create) Email() Email {
	return &email{control: c.control}
}

func (c create) FormBuilder(fields ...*form.FieldBuilder) *form.Builder {
	isCsrfEnabled := c.control.config.Security.Csrf.Enabled
	method := c.control.Request().Method()
	if c.control.Request().Is().Get() {
		method = http.MethodPost
	}
	f := form.New(fields...).
		Method(method).
		Action(c.Link(c.control.route.Name)).
		Request(c.control.request)
	if isCsrfEnabled {
		name := fmt.Sprintf("%s-%s", c.control.route.Name, uniuri.New())
		token := c.control.Csrf().Create(name, c.control.Request().Ip(), c.control.Request().UserAgent())
		f.Csrf(name, token)
	}
	return f
}

func (c create) Link(name string, args ...map[string]any) string {
	if len(name) == 0 {
		return ""
	}
	if c.component != nil {
		if IsFirstCharUpper(name[strings.LastIndex(name, linkLevelDivider)+1:]) {
			return c.createActionLink(name)
		}
	}
	if c.control.Request().Is().Localized() {
		if len(args) == 0 {
			args = make([]map[string]any, 1)
			args[0] = map[string]any{requestVar.Lang: c.control.Request().Lang()}
		}
		if len(args) > 0 {
			args[0][requestVar.Lang] = c.control.Request().Lang()
		}
		localizedRoutes, ok := c.control.core.router.localizedRoutes[c.control.Request().Lang()]
		if !ok {
			return name
		}
		link, ok := c.createLink(localizedRoutes, name, args...)
		if ok {
			return link
		}
	}
	link, _ := c.createLink(c.control.core.router.routes, name, args...)
	return link
}

func (c create) Pdf() Pdf {
	return &pdf{control: c.control}
}

func (c create) Query(qm Map) string {
	n := len(qm)
	if n == 0 {
		return ""
	}
	result := make([]string, n)
	for k, v := range qm {
		result = append(result, fmt.Sprintf("%s=%v", k, v))
	}
	return "?" + strings.Join(result, "&")
}

func (c create) createLink(routes []route.Route, name string, args ...map[string]any) (string, bool) {
	shouldHaveModule := strings.Count(name, linkLevelDivider) == 2
	shouldHaveController := strings.Count(name, linkLevelDivider) >= 1
	shouldAddModulePrefix := len(c.route.Module) > 0 && !shouldHaveModule
	shouldAddControllerPrefix := len(c.route.Controller) > 0 && !shouldHaveController
	for _, rt := range routes {
		if !strings.HasSuffix(rt.Name, name) {
			continue
		}
		if shouldAddControllerPrefix && rt.Controller == c.route.Controller {
			name = rt.Controller + linkLevelDivider + name
		}
		if shouldAddModulePrefix && rt.Module == c.route.Module {
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

func (c create) createActionLink(name string) string {
	name = c.component.Name() + linkLevelDivider + name
	shouldHaveModule := strings.Count(name, linkLevelDivider) == 3
	shouldHaveController := strings.Count(name, linkLevelDivider) >= 2
	shouldAddModulePrefix := len(c.route.Module) > 0 && !shouldHaveModule
	shouldAddControllerPrefix := len(c.route.Controller) > 0 && !shouldHaveController
	if shouldAddControllerPrefix {
		name = c.route.Controller + linkLevelDivider + name
	}
	if shouldAddModulePrefix {
		name = c.route.Module + linkLevelDivider + name
	}
	link := c.Request().Path() + c.Query(Map{queryKey.Action: name})
	if c.control.Config().Component.State == componentState.Cache {
		return link
	}
	params := querystring.New(c.component).
		IgnoreInterface(componentControlInterfaceName).
		Encode()
	if len(params) == 0 {
		return link
	}
	return link + "&" + params
}
