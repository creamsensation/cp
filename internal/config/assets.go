package config

type Assets struct {
	RootPath   string `yaml:"root-path"`
	ConfigPath string `yaml:"config-path"`
	PublicPath string `yaml:"public-path"`
	OutputPath string `yaml:"output-path"`
}
