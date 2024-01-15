package config

type Router struct {
	Localized    bool   `yaml:"localized"`
	PathPrefix   string `yaml:"path-prefix"`
	PreferPrefix bool   `yaml:"prefer-prefix"`
}
