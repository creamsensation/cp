package config

import "time"

type Security struct {
	Csrf      SecurityCsrf                `yaml:"csrf"`
	Firewall  map[string]SecurityFirewall `yaml:"firewall"`
	RateLimit SecurityRateLimit           `yaml:"rate-limit"`
	Role      map[string]SecurityRole     `yaml:"role"`
	Session   SecuritySession             `yaml:"session"`
}

type SecurityCsrf struct {
	Enabled  bool              `yaml:"enabled"`
	Duration time.Duration     `yaml:"duration"`
	Clean    SecurityCsrfClean `yaml:"clean"`
}

type SecurityCsrfClean struct {
	IgnoreRoutes []string `yaml:"ignore-routes"`
}

type SecurityFirewall struct {
	Enabled       bool     `yaml:"enabled"`
	Invert        bool     `yaml:"invert"`
	Modules       []string `yaml:"modules"`
	Controllers   []string `yaml:"controllers"`
	Routes        []string `yaml:"routes"`
	Patterns      []string `yaml:"patterns"`
	RedirectRoute string   `yaml:"redirect_route"`
	Roles         []string `yaml:"roles"`
	Secret        string   `yaml:"secret"`
}

type SecurityRateLimit struct {
	Enabled  bool          `yaml:"enabled"`
	Attempts int           `yaml:"attempts"`
	Interval time.Duration `yaml:"interval"`
}

type SecurityRole struct {
	Super bool `yaml:"super"`
}

type SecuritySession struct {
	Duration time.Duration `yaml:"duration"`
}
