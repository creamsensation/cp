package cp

import (
	"net/http"
	
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/header"
	"github.com/creamsensation/cp/internal/constant/requestVar"
	"github.com/creamsensation/cp/internal/requester"
)

type Request interface {
	ContentType() string
	Form() requester.Form
	Header() http.Header
	Host() string
	Ip() string
	Is() requester.Is
	Lang() string
	Method() string
	Path() string
	Query(key string, defaultValue ...string) string
	Raw() *http.Request
	Route() string
	UserAgent() string
	Var(key string, defaultValue ...string) string
}

type request struct {
	*control
}

func (r request) ContentType() string {
	return r.request.Header.Get(header.ContentType)
}

func (r request) Form() requester.Form {
	return requester.CreateForm(r.request)
}

func (r request) Header() http.Header {
	return r.request.Header
}

func (r request) Host() string {
	return r.Protocol() + "://" + r.request.Host
}

func (r request) Ip() string {
	ip := r.request.Header.Get(header.Ip)
	if len(ip) == 0 {
		return "localhost"
	}
	return ip
}

func (r request) Is() requester.Is {
	return requester.CreateIs(r.request, r.core.router.localizedPathMatcher)
}

func (r request) Lang() string {
	if !r.core.router.localized {
		return ""
	}
	lc := Var[string](r.control, requestVar.Lang)
	if len(lc) == 0 {
		return r.Cookie().Get(cookieName.Lang)
	}
	return lc
}

func (r request) Method() string {
	return r.request.Method
}

func (r request) Path() string {
	return r.request.URL.Path
}

func (r request) Protocol() string {
	if r.request.TLS == nil {
		return "http"
	}
	return "https"
}

func (r request) Query(key string, defaultValue ...string) string {
	return Query[string](r.control, key)
}

func (r request) Raw() *http.Request {
	return r.request
}

func (r request) Route() string {
	return r.control.route.Name
}

func (r request) UserAgent() string {
	return r.request.Header.Get(header.UserAgent)
}

func (r request) Var(key string, defaultValue ...string) string {
	v := Var[string](r.control, key)
	if len(v) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return v
}
