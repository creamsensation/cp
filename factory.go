package cp

import (
	"fmt"
	"net/http"
	
	"github.com/creamsensation/gox"
	"github.com/creamsensation/hx"
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/csrf"
	"github.com/creamsensation/form"
)

type Factory interface {
	Component(ct MandatoryComponent) gox.Node
	Defer(link string, nodes ...gox.Node) gox.Node
	Form(fields ...*form.FieldBuilder) *form.Builder
}

type factory struct {
	ctx *ctx
}

func (f factory) Component(ct MandatoryComponent) gox.Node {
	var action string
	_ = f.ctx.Parse().Query(Action, &action)
	return createComponent(ct, f.ctx, f.ctx.route, action).render()
}

func (f factory) Defer(link string, nodes ...gox.Node) gox.Node {
	return gox.Div(
		hx.Get(link),
		hx.Trigger("load"),
		hx.Swap(hx.SwapOuterHtml),
		hx.Headers(Map{hx.RequestHeaderTrigger: "load"}),
		gox.Fragment(nodes...),
	)
}

func (f factory) Form(fields ...*form.FieldBuilder) *form.Builder {
	isCsrfEnabled := f.ctx.Config().Security.Csrf != nil
	method := f.ctx.Request().Method()
	if f.ctx.Request().Is().Get() {
		method = http.MethodPost
	}
	link := f.ctx.Generate().Link(f.ctx.route.Name)
	if len(f.ctx.Request().Raw().URL.Query()) > 0 {
		q := make(Map)
		for k := range f.ctx.Request().Raw().URL.Query() {
			q[k] = f.ctx.Request().Raw().URL.Query().Get(k)
		}
		link += f.ctx.Generate().Query(q)
	}
	r := form.New(fields...).
		Method(method).
		Action(link).
		Request(f.ctx.r).
		Messages(f.ctx.config.Localization.Form)
	if isCsrfEnabled && !f.ctx.Request().Is().Action() {
		name := fmt.Sprintf("%s-%s", f.ctx.route.Name, uniuri.New())
		token := f.ctx.Csrf().MustCreate(
			csrf.Token{
				Name: name, UserAgent: f.ctx.Request().UserAgent(), Ip: f.ctx.Request().Ip(),
			},
		)
		r.Csrf(name, token)
	}
	return r
}
