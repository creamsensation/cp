package cp

import (
	"fmt"
	
	"github.com/creamsensation/form"
	"github.com/creamsensation/gox"
	"github.com/creamsensation/hx"
)

type HxHandle interface {
	Click(name string) gox.Node
	Submit(name string) gox.Node
}

type hxHandle struct {
	control *control
}

func (h hxHandle) Click(name string) gox.Node {
	token := h.control.Csrf().Create(name, h.control.Request().Ip(), h.control.Request().UserAgent())
	id := "hx-handler-" + name
	h.control.Response().Hx().Update(id)
	return gox.Fragment(
		gox.Id(id),
		hx.Post(h.control.Link(name)),
		hx.Trigger("click"),
		gox.If(
			h.control.config.Security.Csrf.Enabled,
			hx.Vals(fmt.Sprintf(`{"%s":"%s","%s":"%s"}`, form.CsrfName, name, form.CsrfToken, token)),
		),
	)
}

func (h hxHandle) Submit(name string) gox.Node {
	id := "hx-handler-" + name
	h.control.Response().Hx().Update(id)
	return gox.Fragment(
		gox.Id(id),
		hx.Post(h.control.Link(name)),
		hx.Trigger("submit"),
	)
}
