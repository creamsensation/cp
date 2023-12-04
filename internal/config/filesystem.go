package config

type Filesystem struct {
	Driver      string `yaml:"driver"`
	Dir         string `yaml:"dir"`
	StorageName string `yaml:"storage-name"`
	Endpoint    string `yaml:"endpoint"`
	AccessKey   string `yaml:"access-key"`
	SecretKey   string `yaml:"secret-key"`
}
