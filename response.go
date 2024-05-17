package cp

import (
	"github.com/creamsensation/gox"
	
	"github.com/creamsensation/sender"
)

type Response interface {
	sender.ExtendableSend
	Status(statusCode int) Response
	Refresh() error
	Layout(name string) Response
	Render(nodes ...gox.Node) error
	Intercept() Intercept
}

type response struct {
	*sender.Sender
	ctx    *ctx
	layout *layout
	l      layoutFactory
}

func (r *response) Refresh() error {
	return r.Redirect(r.ctx.Request().Path())
}

func (r *response) Layout(name string) Response {
	if r.layout == nil {
		return r
	}
	l, ok := r.layout.factories[name]
	if !ok {
		panic(ErrorInvalidLayout)
	}
	r.l = l
	return r
}

func (r *response) Status(statusCode int) Response {
	r.StatusCode = statusCode
	return r
}

func (r *response) Render(nodes ...gox.Node) error {
	if r.layout != nil && r.l != nil && !r.ctx.Request().Is().Hx() {
		return r.Html(gox.Render(r.l(r.ctx, nodes...)))
	}
	return r.Html(gox.Render(nodes...))
}

func (r *response) Intercept() Intercept {
	return interceptor{
		Sender: r.Sender,
		err:    r.ctx.err,
	}
}
