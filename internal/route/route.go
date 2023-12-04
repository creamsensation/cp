package route

import (
	"regexp"
	
	"github.com/creamsensation/cp/internal/firewall"
	"github.com/creamsensation/cp/internal/handler"
)

type Route struct {
	Controller       string
	Firewalls        map[string]*firewall.Route
	Fn               handler.Fn
	Layout           string
	Matcher          *regexp.Regexp
	Method           []string
	Module           string
	Name             string
	Ok               bool
	Path             string
	VarsPlaceholders map[string]string
}

func (r Route) Map() map[string]any {
	return map[string]any{
		"Controller": r.Controller,
		"Layout":     r.Layout,
		"Module":     r.Module,
		"Name":       r.Name,
		"Path":       r.Path,
	}
}
