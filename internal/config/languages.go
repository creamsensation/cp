package config

type Languages map[string]Language

type Language struct {
	Enabled bool `yaml:"enabled"`
	Default bool `yaml:"default"`
}
