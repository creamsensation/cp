package cp

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
)

func TestAssetsReader(t *testing.T) {
	rootDir := t.TempDir()
	cfg := config.Assets{
		RootPath:   "assets",
		ConfigPath: "config",
		PublicPath: ".static",
		OutputPath: ".dist",
	}
	stylePath, scriptPath := ".static/.dist/main.css", ".static/.dist/main.js"
	outputDir := fmt.Sprintf("%s/%s/%s/%s", rootDir, cfg.RootPath, cfg.PublicPath, cfg.OutputPath)
	assert.Nil(t, os.MkdirAll(outputDir, os.ModePerm))
	ar := createAssetsReader(rootDir, cfg)
	manifestBts, err := json.Marshal(
		map[string]string{
			stylePath:  stylePath,
			scriptPath: scriptPath,
		},
	)
	assert.Nil(t, err)
	assert.Nil(t, os.WriteFile(fmt.Sprintf("%s/manifest.json", outputDir), manifestBts, os.ModePerm))
	ar.read()
	assert.Equal(t, []string{stylePath}, ar.styles)
	assert.Equal(t, []string{scriptPath}, ar.scripts)
}
