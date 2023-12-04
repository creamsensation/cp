package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	
	"gopkg.in/yaml.v3"
	
	"github.com/creamsensation/cp/env"
)

type Config struct {
	App        App        `yaml:"app"`
	Assets     Assets     `yaml:"assets"`
	Cache      Cache      `yaml:"cache"`
	Component  Component  `yaml:"component"`
	Database   Databases  `yaml:"database"`
	Filesystem Filesystem `yaml:"filesystem"`
	Languages  Languages  `yaml:"languages"`
	Router     Router     `yaml:"router"`
	Security   Security   `yaml:"security"`
	Smtp       Smtp       `yaml:"smtp"`
}

func Parse(dir string) Config {
	var result Config
	if strings.HasPrefix(dir, "/") {
		dir = strings.TrimPrefix(dir, "/")
	}
	dir = fmt.Sprintf("%s/%s", dir, env.Get())
	parsePart(dir, "app", &result)
	parsePart(dir, "assets", &result)
	parsePart(dir, "cache", &result)
	parsePart(dir, "component", &result)
	parsePart(dir, "database", &result)
	parsePart(dir, "filesystem", &result)
	parsePart(dir, "languages", &result)
	parsePart(dir, "router", &result)
	parsePart(dir, "security", &result)
	parsePart(dir, "smtp", &result)
	if !strings.HasPrefix(result.Assets.PublicDir, "/") {
		result.Assets.PublicDir = "/" + result.Assets.PublicDir
	}
	result.Router.PathPrefix = strings.TrimPrefix(strings.TrimSuffix(result.Router.PathPrefix, "/"), "")
	return result
}

func parsePart(dir string, name string, data any) {
	path := fmt.Sprintf("%s/%s.yaml", dir, name)
	configBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if err := yaml.Unmarshal(configBytes, data); err != nil {
		log.Fatalln(err)
	}
}
