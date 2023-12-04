package firewall

type Route struct {
	Enabled       bool
	Invert        bool
	Roles         []string
	RedirectRoute string
	Secret        string
}
