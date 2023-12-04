package config

type Cache struct {
	Adapter  string `yaml:"adapter"`
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}
