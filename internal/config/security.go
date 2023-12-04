package config

import "time"

type Security struct {
	Csrf struct {
		Enabled  bool          `yaml:"enabled"`
		Duration time.Duration `yaml:"duration"`
	} `yaml:"csrf"`
	Firewall  map[string]securityFirewall `yaml:"firewall"`
	RateLimit securityRateLimit           `yaml:"rate-limit"`
	Role      map[string]securityRole     `yaml:"role"`
	Session   securitySession             `yaml:"session"`
}

type securityFirewall struct {
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

type securityRateLimit struct {
	Enabled  bool          `yaml:"enabled"`
	Attempts int           `yaml:"attempts"`
	Interval time.Duration `yaml:"interval"`
}

type securityRole struct {
	Super bool `yaml:"super"`
}

type securitySession struct {
	Duration time.Duration `yaml:"duration"`
}
