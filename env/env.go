package env

import (
	"log"
	"os"
)

const (
	envVar      = "APP_ENV"
	development = "development"
	production  = "production"
)

func EnableProductionMode() {
	if err := os.Setenv(envVar, production); err != nil {
		log.Fatalln(err)
	}
}

func EnableDevelopmentMode() {
	if err := os.Setenv(envVar, development); err != nil {
		log.Fatalln(err)
	}
}

func Get() string {
	return os.Getenv(envVar)
}

func Development() bool {
	return Get() == development
}

func Production() bool {
	return Get() == production
}
