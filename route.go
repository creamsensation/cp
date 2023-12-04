package cp

import (
	"net/http"
	
	"github.com/creamsensation/cp/internal/handler"
	"github.com/creamsensation/cp/internal/route"
)

type Routes []*route.Builder

func Route(configs ...route.Config) *route.Builder {
	return route.CreateBuilder(configs...)
}

func Name(name string) route.Config {
	return route.CreateConfig(route.ConfigName, name)
}

func Path(path route.PathValue) route.Config {
	return route.CreatePathConfig(path)
}

func Layout(name string) route.Config {
	return route.CreateConfig(route.ConfigLayout, name)
}

func Localize() route.Config {
	return route.CreateConfig(route.ConfigLocalize, true)
}

func Handler(fn handler.Fn) route.Config {
	return route.CreateConfig(route.ConfigHandler, fn)
}

func Method(methods ...route.Config) route.Config {
	return route.CreateMethodConfig(methods...)
}

func Get(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodGet, path...)
}

func Post(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodPost, path...)
}

func Put(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodPut, path...)
}

func Patch(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodPatch, path...)
}

func Delete(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodDelete, path...)
}

func Options(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodOptions, path...)
}

func Head(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodHead, path...)
}

func Trace(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodTrace, path...)
}

func Connect(path ...route.PathValue) route.Config {
	return route.CreatePath(http.MethodConnect, path...)
}
