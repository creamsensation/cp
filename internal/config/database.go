package config

type Databases = map[string]Database

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"db"`
	Ssl      string `yaml:"sslmode"`
	CertPath string `yaml:"certpath"`
}
