package cp

import (
	"fmt"
	"net/http"
	
	"github.com/creamsensation/util/constant/contentType"
	"github.com/creamsensation/util/constant/dataType"
	"github.com/creamsensation/util/constant/header"
)

type Handler func(c Ctx) error

type handler struct {
	core   *core
	method string
	path   string
	name   string
}

func (h handler) create(fn Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		matchedRoute := h.matchRoute(r.URL.Path)
		c := createContext(
			ctxParam{
				assets:       h.core.assets,
				config:       h.core.router.config,
				layout:       h.core.layout,
				r:            r,
				w:            w,
				matchedRoute: matchedRoute,
				routes:       h.core.router.routes,
			},
		)
		if h.core.router.config.Router.Recover {
			defer h.createRecover(c)
		}
		for _, middleware := range h.applyInternalMiddlewares(matchedRoute, h.core.router.middlewares) {
			c.mu.Lock()
			c.err = middleware(c)
			if c.err != nil {
				c.mu.Unlock()
				h.createResponse(c)
				return
			}
		}
		if len(c.response.DataType) == 0 {
			err := fn(c)
			if err != nil {
				c.err = err
			}
		}
		h.createResponse(c)
	}
}

func (h handler) applyInternalMiddlewares(matchedRoute *Route, middlewares []Handler) []Handler {
	r := make([]Handler, 0)
	if h.core.config.Localization.Enabled {
		r = append(r, createLangMiddleware())
	}
	if h.core.config.Security.Csrf != nil {
		r = append(r, createCsrfMiddleware())
	}
	if len(matchedRoute.Firewalls) > 0 {
		r = append(r, createFirewallMiddleware(matchedRoute.Firewalls))
	}
	r = append(r, middlewares...)
	return r
}

func (h handler) createResponse(c *ctx) {
	if c.err != nil {
		c.w.Header().Set(header.ContentType, contentType.Text)
		if c.response.StatusCode == http.StatusOK {
			c.w.WriteHeader(http.StatusInternalServerError)
		}
		if c.response.StatusCode != http.StatusOK {
			c.w.WriteHeader(c.response.StatusCode)
		}
		_, err := c.w.Write([]byte(c.err.Error()))
		if err != nil {
			http.Error(c.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	if c.response.DataType == dataType.Redirect {
		if c.response.StatusCode == http.StatusOK {
			c.response.StatusCode = http.StatusFound
		}
		http.Redirect(c.w, c.r, c.response.Value, c.response.StatusCode)
		return
	}
	if c.response.DataType == dataType.Stream {
		c.w.Header().Set(header.ContentDisposition, fmt.Sprintf("attachment;filename=%s", c.response.Value))
		c.w.Header().Set(header.ContentLength, fmt.Sprintf("%d", len(c.response.Bytes)))
	}
	c.w.Header().Set(header.ContentType, c.response.ContentType)
	c.w.WriteHeader(c.response.StatusCode)
	if _, c.err = c.w.Write(c.response.Bytes); c.err != nil {
		http.Error(c.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h handler) createRecover(c *ctx) {
	if e := recover(); e != nil {
		err, ok := e.(error)
		if !ok {
			http.Error(c.w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if c.response.StatusCode == http.StatusOK || c.response.StatusCode == http.StatusBadRequest {
			c.response.StatusCode = http.StatusInternalServerError
		}
		c.err = err
		err = h.core.errorHandler(c)
		c.err = nil
		if err != nil {
			c.err = err
		}
		h.createResponse(c)
	}
}

func (h handler) matchRoute(path string) *Route {
	for _, r := range *h.core.router.routes {
		if r.Matcher.MatchString(path) {
			return r
		}
	}
	return nil
}
