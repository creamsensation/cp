package firewall

import "regexp"

type Firewall struct {
	Enabled       bool
	Invert        bool
	Modules       []string
	Controllers   []string
	Routes        []string
	Patterns      []string
	Matchers      []*regexp.Regexp
	RedirectRoute string
	Roles         []string
	Secret        string
}
