package config

import (
	"fmt"
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/env"
)

func TestConfig(t *testing.T) {
	env.EnableDevelopmentMode()
	rootDir := t.TempDir()
	envDir := fmt.Sprintf("%s/%s", rootDir, env.Get())
	path := fmt.Sprintf("%s/app.yaml", envDir)
	assert.Nil(t, os.MkdirAll(envDir, os.ModePerm))
	assert.Nil(
		t,
		os.WriteFile(
			path, []byte(`app:
  name: "test"
  port: 8000`),
			os.ModePerm,
		),
	)
	t.Run(
		"parse", func(t *testing.T) {
			cfg := Parse(rootDir)
			assert.Equal(t, "test", cfg.App.Name)
			assert.Equal(t, 8000, cfg.App.Port)
		},
	)
}
