package cp

import (
	"net/http"
	"net/url"
	
	"github.com/creamsensation/hx"
	
	"github.com/creamsensation/util/constant/header"
)

type Request interface {
	ContentType() string
	Form() url.Values
	Header() http.Header
	Host() string
	Ip() string
	Is() RequestIs
	Method() string
	Name() string
	Origin() string
	Path() string
	PathValue(key string, defaultValue ...string) string
	QueryParam(key string, defaultValue ...string) string
	Protocol() string
	Raw() *http.Request
	UserAgent() string
}

type RequestIs interface {
	Get() bool
	Post() bool
	Put() bool
	Patch() bool
	Delete() bool
	Action() bool
	Hx() bool
	Options() bool
	Head() bool
	Connect() bool
	Trace() bool
}

type request struct {
	r     *http.Request
	route *Route
}

func (r request) ContentType() string {
	return r.r.Header.Get(header.ContentType)
}

func (r request) Form() url.Values {
	return r.r.Form
}

func (r request) Header() http.Header {
	return r.r.Header
}

func (r request) Host() string {
	return r.Protocol() + "://" + r.r.Host
}

func (r request) Ip() string {
	return r.r.Header.Get("X-Forwarded-For")
}

func (r request) Is() RequestIs {
	return r
}

func (r request) Method() string {
	return r.r.Method
}

func (r request) Name() string {
	return r.route.Name
}

func (r request) Origin() string {
	return r.r.Header.Get(header.Origin)
}

func (r request) Path() string {
	return r.r.URL.Path
}

func (r request) PathValue(key string, defaultValue ...string) string {
	value := r.r.PathValue(key)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (r request) QueryParam(key string, defaultValue ...string) string {
	value := r.r.URL.Query().Get(key)
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (r request) Protocol() string {
	if r.r.TLS == nil {
		return "http"
	}
	return "https"
}

func (r request) Raw() *http.Request {
	return r.r
}

func (r request) UserAgent() string {
	return r.r.Header.Get(header.UserAgent)
}

func (r request) Get() bool {
	return r.r.Method == http.MethodGet
}

func (r request) Post() bool {
	return r.r.Method == http.MethodPost
}

func (r request) Put() bool {
	return r.r.Method == http.MethodPut
}

func (r request) Patch() bool {
	return r.r.Method == http.MethodPatch
}

func (r request) Delete() bool {
	return r.r.Method == http.MethodDelete
}

func (r request) Action() bool {
	return len(r.r.URL.Query().Get(Action)) > 0
}

func (r request) Hx() bool {
	return r.r.Header.Get(hx.RequestHeaderRequest) == "true"
}

func (r request) Options() bool {
	return r.r.Method == http.MethodOptions
}

func (r request) Head() bool {
	return r.r.Method == http.MethodHead
}

func (r request) Connect() bool {
	return r.r.Method == http.MethodConnect
}

func (r request) Trace() bool {
	return r.r.Method == http.MethodTrace
}
