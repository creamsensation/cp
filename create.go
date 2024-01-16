package cp

import (
	"errors"
	"fmt"
	"net/http"
	
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/form"
	"github.com/creamsensation/gox"
	"github.com/creamsensation/hx"
)

type Create interface {
	Component(component component) gox.Node
	Defer(link string, nodes ...gox.Node) gox.Node
	FormBuilder(fields ...*form.FieldBuilder) *form.Builder
}

type create struct {
	*control
	component component
}

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
		hx.Headers(Map{hx.RequestHeaderTrigger: "load"}),
		gox.Fragment(nodes...),
	)
}

func (c create) FormBuilder(fields ...*form.FieldBuilder) *form.Builder {
	isCsrfEnabled := c.control.config.Security.Csrf.Enabled
	method := c.control.Request().Method()
	if c.control.Request().Is().Get() {
		method = http.MethodPost
	}
	f := form.New(fields...).
		Method(method).
		Action(c.Generate().Link().Name(c.control.route.Name)).
		Request(c.control.request)
	if len(c.core.form.errors) > 0 {
		errs := make(map[string]error)
		for k, e := range c.core.form.errors {
			errs[k] = errors.New(c.Translate(e))
		}
		f.Errors(errs)
	}
	if isCsrfEnabled {
		name := fmt.Sprintf("%s-%s", c.control.route.Name, uniuri.New())
		token := c.control.Csrf().Create(name, c.control.Request().Ip(), c.control.Request().UserAgent())
		f.Csrf(name, token)
	}
	return f
}
