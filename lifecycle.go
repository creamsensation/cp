package cp

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/expiration"
	"github.com/creamsensation/cp/internal/constant/header"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/cp/internal/constant/queryKey"
	"github.com/creamsensation/cp/internal/constant/requestVar"
	"github.com/creamsensation/cp/internal/dev"
	"github.com/creamsensation/cp/internal/firewall"
	"github.com/creamsensation/cp/internal/handler"
	"github.com/creamsensation/cp/internal/result"
	"github.com/creamsensation/cp/internal/route"
)

type lifecycle struct {
	control  *control
	core     *core
	request  *http.Request
	response http.ResponseWriter
	written  bool
}

func createLifecycle(control *control) *lifecycle {
	return &lifecycle{
		control:  control,
		core:     control.core,
		request:  control.request,
		response: control.response,
	}
}

func (l *lifecycle) run() {
	defer l.recoverErrors()
	l.runRouteMatcher()
	l.extractVars()
	if ok := l.validateCanonization(); !ok {
		return
	}
	if ok := l.runMiddlewares(); !ok {
		return
	}
	if ok := l.runFirewall(); !ok {
		return
	}
	l.runResponseResultCreator()
}

func (l *lifecycle) runMiddlewares() bool {
	for _, handlerFn := range l.core.router.middlewares {
		l.callFn(handlerFn)
		if len(l.control.result.Content) == 0 {
			continue
		}
		l.writeResult()
		return false
	}
	return true
}

func (l *lifecycle) runRouteMatcher() {
	var match route.Route
	var isPathLocalized bool
	if l.core.router.localizedPathMatcher != nil {
		isPathLocalized = l.control.Request().Is().Localized()
	}
	if isPathLocalized {
		for _, routes := range l.core.router.localizedRoutes {
			for _, r := range routes {
				if r.Matcher.MatchString(l.control.Request().Path()) && slices.Contains(
					r.Method, l.control.Request().Method(),
				) {
					match = r
					break
				}
			}
		}
	}
	if !isPathLocalized {
		for _, r := range l.core.router.routes {
			if r.Matcher == nil {
				continue
			}
			if r.Matcher.MatchString(l.control.Request().Path()) && slices.Contains(r.Method, l.control.Request().Method()) {
				match = r
				break
			}
		}
	}
	if l.core.router.localized && !match.Ok && !isPathLocalized {
		var defaultPath string
		isPrefix := len(l.core.config.Router.PathPrefix) > 0
		if !isPrefix {
			defaultPath = fmt.Sprintf("/%s", l.core.router.defaultLanguage)
		}
		if isPrefix && !l.core.config.Router.PreferPrefix {
			defaultPath = fmt.Sprintf("/%s/%s", l.core.router.defaultLanguage, l.core.config.Router.PathPrefix)
		}
		if isPrefix && l.core.config.Router.PreferPrefix {
			defaultPath = fmt.Sprintf("/%s/%s", l.core.config.Router.PathPrefix, l.core.router.defaultLanguage)
		}
		if l.control.Request().Path() != defaultPath {
			l.setHttpRedirect(defaultPath, http.StatusFound)
			return
		}
	}
	if !match.Ok {
		l.createErrorResult(
			http.StatusNotFound, errors.New(fmt.Sprintf("route [%s] not found", l.control.Request().Path())),
		)
		return
	}
	l.control.route = l.validateLayout(match)
}

func (l *lifecycle) runFirewall() bool {
	if l.control.route.Firewalls == nil {
		return true
	}
	if len(l.control.route.Firewalls) == 0 {
		return true
	}
	firewalls := make(map[string]*firewall.Route)
	for name, f := range l.control.route.Firewalls {
		if !f.Enabled {
			continue
		}
		firewalls[name] = f
	}

	if len(firewalls) == 0 {
		return true
	}
	isAuth := l.control.Auth().Session().Exists()
	for _, f := range firewalls {
		if f.Invert && isAuth {
			l.setHttpRedirect(l.control.Generate().Link().Name(f.RedirectRoute), http.StatusFound)
			return false
		}
	}
	ss := getSession(l.control)
	if ss.Super {
		return true
	}
	validFirewalls := make([]bool, 0)
	secret := l.control.Request().Header().Get(header.Secret)
	for _, f := range firewalls {
		var valid bool
		if len(f.Secret) > 0 && f.Secret == secret {
			valid = true
		}
		if f.Invert || (!f.Invert && isAuth && len(f.Roles) == 0) {
			valid = true
		}
		for _, roleName := range f.Roles {
			if slices.Contains(ss.Roles, roleName) {
				valid = true
			}
		}
		if !valid && len(f.RedirectRoute) > 0 {
			l.setHttpRedirect(l.control.Generate().Link().Name(f.RedirectRoute), http.StatusFound)
			return false
		}
		if !valid && len(f.RedirectRoute) == 0 {
			l.createErrorResult(http.StatusForbidden, ErrorForbidden)
			return false
		}
		validFirewalls = append(validFirewalls, valid)
	}
	if len(validFirewalls) == len(firewalls) {
		return true
	}
	l.createErrorResult(http.StatusForbidden, ErrorForbidden)
	return false
}

func (l *lifecycle) runResponseResultCreator() {
	if !l.control.route.Ok {
		return
	}
	l.prepareDev()
	l.createSecurityHeaders()
	l.validateLanguage()
	l.callFn(l.control.route.Fn)
	l.writeResult()
}

func (l *lifecycle) validateLanguage() {
	var lc string
	languagesExist := l.core.languagesExist()
	if languagesExist && l.control.core.router.localized {
		lc = l.control.Request().Var(requestVar.Lang)
		if len(lc) == 0 {
			for k, v := range l.control.config.Languages {
				if v.Enabled && v.Default {
					lc = k
					break
				}
			}
		}
	}
	if languagesExist && !l.control.core.router.localized {
		lc = l.control.Cookie().Get(cookieName.Lang)
		if len(lc) == 0 {
			for k, v := range l.control.config.Languages {
				if v.Enabled && v.Default {
					lc = k
					break
				}
			}
		}
		lc = l.control.Request().Query(queryKey.Lang, lc)
	}
	if languagesExist {
		l.control.Cookie().Set(cookieName.Lang, lc, expiration.Lang)
	}
}

func (l *lifecycle) createSecurityHeaders() {
	l.response.Header().Set("X-Frame-Options", "SAMEORIGIN")
	l.response.Header().Set("X-XSS-Protection", "1; mode=block")
	l.response.Header().Set("X-Content-Type-Options", "nosniff")
	l.response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	l.response.Header().Set("Cross-Origin-Resource-Policy", "same-site")
	l.response.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	l.response.Header().Set("Vary", "origin")
	l.response.Header().Set(
		"Cross-Origin-Embedder-Policy", "require-corp",
	) // Pokud necemu vyslovene nedas crossorigin, brani v loadu
	l.response.Header().Set("Access-Control-Allow-Origin", l.control.Request().Host())
}

func (l *lifecycle) writeResult() {
	switch l.control.result.ResultType {
	case result.Error, result.Json, result.Render, result.Text:
		l.createResult()
	case result.Stream:
		l.createStreamResult()
	case result.Redirect:
		l.createRedirectResult()
	}
}

func (l *lifecycle) callFn(fn handler.Fn) {
	ref := reflect.ValueOf(fn)
	if !ref.IsValid() {
		l.createErrorResult(http.StatusInternalServerError, errors.New("invalid handler function"))
		return
	}
	args := make([]reflect.Value, 0)
	for i := 0; i < ref.Type().NumIn(); i++ {
		in := ref.Type().In(i)
		if in.String() == controlInterfaceName {
			args = append(args, reflect.ValueOf(l.control))
			continue
		}
		dep := provide[any](l.control, in, in.String())
		if dep == nil {
			l.createErrorResult(
				http.StatusInternalServerError, errors.New(fmt.Sprintf("dependency [%s] does not exist", in.String())),
			)
			return
		}
		depRef := reflect.ValueOf(dep)
		if !depRef.IsValid() {
			l.createErrorResult(
				http.StatusInternalServerError, errors.New(fmt.Sprintf("dependency [%s] is invalid", in.String())),
			)
			return
		}
		args = append(args, depRef)
	}
	r := ref.Call(args)[0].Interface()
	if r == nil {
		return
	}
	if l.control.result.StatusCode > 0 && len(l.control.result.Content) > 0 {
		return
	}
	l.control.result = r.(result.Result)
}

func (l *lifecycle) createResult() {
	l.setHttpResult()
}

func (l *lifecycle) createStreamResult() {
	l.setHttpStreamResult()
}

func (l *lifecycle) createRedirectResult() {
	l.setHttpRedirect(l.control.result.Content, l.control.result.StatusCode)
}

func (l *lifecycle) createErrorResult(statusCode int, err any) {
	e := errors.New(fmt.Sprintf("%v", err))
	l.control.result = result.CreateError(createErrorPage(l.control, statusCode, e), statusCode, e)
	panic(e)
}

func (l *lifecycle) validateCanonization() bool {
	if !l.control.route.Ok {
		return true
	}
	path := l.control.route.Path
	for k, p := range l.control.route.VarsPlaceholders {
		path = strings.Replace(path, p, Var[string](l.control, k), 1)
	}
	if len(path) == 0 {
		return true
	}
	encodedQuery := l.control.Request().Raw().URL.Query().Encode()
	if l.control.Request().Path() != path && (len(encodedQuery) == 0 || encodedQuery == "?") {
		l.setHttpRedirect(path, http.StatusMovedPermanently)
		return false
	}
	return true
}

func (l *lifecycle) extractVars() {
	if len(l.control.route.Path) == 0 {
		return
	}
	pathParts := strings.Split(l.control.route.Path, "/")
	for i, item := range strings.Split(l.control.Request().Path(), "/") {
		if len(item) == 0 {
			continue
		}
		if strings.HasPrefix(pathParts[i], "{") && strings.HasSuffix(pathParts[i], "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(pathParts[i], "{"), "}")
			if strings.Contains(key, ":") {
				key = key[:strings.Index(key, ":")]
			}
			l.control.vars[key] = item
		}
	}
}

func (l *lifecycle) recoverErrors() {
	e := recover()
	if e == nil {
		return
	}
	if *l.control.statusCode == 0 {
		*l.control.statusCode = http.StatusInternalServerError
	}
	if len(l.control.result.ContentType) == 0 {
		err := errors.New(fmt.Sprintf("%+v", e))
		l.control.result = result.CreateError(
			createErrorPage(l.control, *l.control.statusCode, err), *l.control.statusCode, err,
		)
	}
	if env.Development() && !l.written {
		l.response.Header().Set(header.ContentType, l.control.result.ContentType)
		if _, err := l.response.Write([]byte(l.control.result.Content)); err != nil {
			l.setHttpError(fmt.Sprintf("%+v", e), l.control.result.StatusCode)
		}
		l.written = true
	}
	if !env.Development() && e != nil {
		l.setHttpError(fmt.Sprintf("%+v", e), l.control.result.StatusCode)
	}
	if !env.Development() && e == nil {
		l.setHttpError(http.StatusText(l.control.result.StatusCode), l.control.result.StatusCode)
	}
}

func (l *lifecycle) validateLayout(r route.Route) route.Route {
	if r.Ok && len(r.Layout) == 0 {
		r.Layout = naming.Main
	}
	return r
}

func (l *lifecycle) setHttpResult() {
	if l.written {
		return
	}
	if *l.control.statusCode == 0 {
		*l.control.statusCode = http.StatusOK
	}
	l.response.Header().Set(header.ContentType, l.control.result.ContentType)
	l.response.WriteHeader(*l.control.statusCode)
	if _, err := l.response.Write([]byte(l.control.result.Content)); err != nil {
		l.createErrorResult(http.StatusInternalServerError, errors.New("response data could not be written"))
	}
	l.written = true
}

func (l *lifecycle) setHttpStreamResult() {
	if l.written {
		return
	}
	if *l.control.statusCode == 0 {
		*l.control.statusCode = http.StatusOK
	}
	l.response.Header().Set(header.ContentType, l.control.result.ContentType)
	l.response.Header().Set(header.ContentDisposition, fmt.Sprintf("attachment;filename=%s", l.control.result.Content))
	l.response.Header().Set(header.ContentLength, fmt.Sprintf("%d", len(l.control.result.Data)))
	l.response.WriteHeader(*l.control.statusCode)
	if _, err := l.response.Write(l.control.result.Data); err != nil {
		l.createErrorResult(http.StatusInternalServerError, errors.New("response data could not be written"))
	}
	l.written = true
}

func (l *lifecycle) setHttpError(msg string, code int) {
	if l.written {
		return
	}
	http.Error(l.response, msg, code)
	l.written = true
}

func (l *lifecycle) setHttpRedirect(url string, code int) {
	if l.written {
		return
	}
	http.Redirect(l.response, l.request, url, code)
	l.written = true
}

func (l *lifecycle) prepareDev() {
	if !env.Development() {
		return
	}
	r := l.control.route
	r.Path = l.control.Request().Path()
	l.control.dev.(dev.Internal).SetRoute(r)
}
