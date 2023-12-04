package config

type Router struct {
	PathPrefix   string `yaml:"path-prefix"`
	PreferPrefix bool   `yaml:"prefer-prefix"`
}
