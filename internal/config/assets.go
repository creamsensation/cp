package config

type Assets struct {
	PublicDir  string `yaml:"public-dir"`
	EntryPath  string `yaml:"entry-path"`
	PublicPath string `yaml:"public-path"`
	ConfigPath string `yaml:"config-path"`
}
