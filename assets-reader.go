package cp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	
	"github.com/creamsensation/cp/internal/config"
)

type assetsReader struct {
	config       config.Assets
	manifestPath string
	rootDir      string
	scripts      []string
	styles       []string
	timestamp    time.Time
}

const (
	sourcemapSuffix = ".map"
)

func createAssetsReader(dir string, config config.Assets) *assetsReader {
	ar := &assetsReader{
		config:  config,
		scripts: make([]string, 0),
		styles:  make([]string, 0),
	}
	ar.rootDir = ar.createRootDir(dir)
	ar.manifestPath = fmt.Sprintf(
		"%s/%s/%s/manifest.json", ar.rootDir, config.PublicPath, config.OutputPath,
	)
	ar.createConfig()
	return ar
}

func (a *assetsReader) createRootDir(wd string) string {
	path := []string{wd}
	if a.config.RootPath != "." {
		path = append(path, a.config.RootPath)
	}
	return strings.Join(path, "/")
}

func (a *assetsReader) createConfig() {
	path := fmt.Sprintf(`%s/assets.config.ts`, a.rootDir)
	if err := os.MkdirAll(a.rootDir, os.ModePerm); err != nil {
		fmt.Printf("%s dir not exists\n", a.rootDir)
		return
	}
	if err := os.WriteFile(
		path, []byte(fmt.Sprintf(
			`export const assetsConfig = {
	outputPath: '%s',
	publicPath: '%s',
	rootPath: '%s',
}
`, a.config.OutputPath, a.config.PublicPath, a.config.RootPath,
		)), os.ModePerm,
	); err != nil {
		fmt.Println(err)
	}
}

func (a *assetsReader) read() {
	stats, err := os.Stat(a.manifestPath)
	if os.IsNotExist(err) || !os.IsNotExist(err) && stats.ModTime().Compare(a.timestamp) == 0 {
		return
	}
	a.styles = make([]string, 0)
	a.scripts = make([]string, 0)
	a.timestamp = stats.ModTime()
	mb, err := os.ReadFile(a.manifestPath)
	if err != nil {
		log.Fatalln(err)
	}
	manifest := make(map[string]string)
	if err := json.Unmarshal(mb, &manifest); err != nil {
		log.Fatalln(err)
	}
	for _, path := range manifest {
		if strings.HasSuffix(path, sourcemapSuffix) {
			continue
		}
		if strings.HasSuffix(path, ".css") {
			a.styles = append(a.styles, strings.TrimPrefix(path, a.config.RootPath))
		}
		if strings.HasSuffix(path, ".js") {
			a.scripts = append(a.scripts, strings.TrimPrefix(path, a.config.RootPath))
		}
	}
}
