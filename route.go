package cp

import (
	"regexp"
	
	"github.com/creamsensation/firewall"
)

type RouteConfig struct {
	Type  int
	Value any
}

type Route struct {
	Lang      string
	Path      string
	Name      string
	Matcher   *regexp.Regexp
	Methods   []string
	Firewalls []firewall.Firewall
}

const (
	routeMethod = iota
	routeName
)

func Method(method ...string) RouteConfig {
	return RouteConfig{
		Type:  routeMethod,
		Value: method,
	}
}

func Name(name string) RouteConfig {
	return RouteConfig{
		Type:  routeName,
		Value: name,
	}
}
