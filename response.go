package cp

import (
	"encoding/json"
	"net/http"
	"strings"
	
	"github.com/creamsensation/gox"
	htmx "github.com/creamsensation/hx"
	
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/responder/hx"
	"github.com/creamsensation/cp/internal/result"
)

type Result interface{}

type Response interface {
	Error(err error) Result
	File(name string, data []byte) Result
	Json(value any) Result
	Header() http.Header
	Hx() hx.Response
	Raw() http.ResponseWriter
	Redirect(name string) Result
	Refresh() Result
	Render(nodes ...gox.Node) Result
	Status(statusCode int) Response
	Text(text string) Result
	Empty() Result
}

type response struct {
	control *control
}

func (r response) Error(err error) Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusInternalServerError
	}
	return result.CreateError(createErrorPage(r.control, *r.control.statusCode, err), *r.control.statusCode, err)
}

func (r response) File(name string, data []byte) Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusOK
	}
	return result.CreateStream(name, data, *r.control.statusCode)
}

func (r response) Json(value any) Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusOK
	}
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return result.CreateError(
			createErrorPage(r.control, *r.control.statusCode, err), http.StatusInternalServerError, err,
		)
	}
	return result.CreateJson(string(valueBytes), *r.control.statusCode)
}

func (r response) Header() http.Header {
	return r.control.response.Header()
}

func (r response) Hx() hx.Response {
	return r.control.hxResponse
}

func (r response) Raw() http.ResponseWriter {
	return r.control.response
}

func (r response) Redirect(name string) Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusFound
	}
	r.control.flash.store()
	if strings.HasPrefix(name, "/") {
		return result.CreateRedirect(name, *r.control.statusCode)
	}
	return result.CreateRedirect(r.control.Link(name), *r.control.statusCode)
}

func (r response) Refresh() Result {
	return r.Redirect(r.control.route.Name)
}

func (r response) Render(nodes ...gox.Node) Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusOK
	}
	if r.shouldUseHx() {
		if env.Development() {
			nodes = append(nodes, r.control.Dev().Tool())
		}
		for i, n := range nodes {
			if !r.control.hxResponse.Exists(gox.GetAttribute[string](n, "id")) {
				continue
			}
			nodes[i] = gox.Append(n, htmx.SwapOob())
		}
		r.control.hxResponse.PrepareHeaders()
		return result.CreateHtml(gox.Render(nodes...), *r.control.statusCode)
	}
	return result.CreateHtml(gox.Render(r.getLayout()(r.control, nodes...)), *r.control.statusCode)
}

func (r response) Status(statusCode int) Response {
	*r.control.statusCode = statusCode
	return r
}

func (r response) Text(text string) Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusOK
	}
	return result.CreateText(text, *r.control.statusCode)
}

func (r response) Empty() Result {
	if r.control.main != nil || *r.control.statusCode == 0 {
		*r.control.statusCode = http.StatusOK
	}
	return result.CreateText("", *r.control.statusCode)
}

func (r response) shouldUseHx() bool {
	return r.control.Request().Is().Hx()
}

func (r response) getLayout() uiLayoutFn {
	l, ok := r.control.core.ui.layouts[r.control.route.Layout]
	if !ok {
		return func(c Control, nodes ...gox.Node) gox.Node {
			return gox.Fragment(nodes...)
		}
	}
	return l
}
