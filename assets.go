package cp

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"
)

type assets struct {
	dir     string
	public  string
	styles  []string
	scripts []string
}

type viteManifest struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	Css     []string `json:"css"`
}

const (
	manifestFilaname = "manifest.json"
)

func (a *assets) getDistDir() string {
	if len(a.dir) == 0 || len(a.public) == 0 || !strings.Contains(a.dir, a.public) {
		return ""
	}
	return a.dir[strings.Index(a.dir, a.public)+len(a.public):]
}

func (a *assets) read() error {
	distDir := a.getDistDir()
	if len(distDir) == 0 {
		return nil
	}
	filePath, err := url.JoinPath(a.dir, manifestFilaname)
	if err != nil {
		return err
	}
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	assetsMap := make(map[string]viteManifest)
	if err := json.Unmarshal(bytes, &assetsMap); err != nil {
		return err
	}
	for _, item := range assetsMap {
		var err error
		item.File, err = url.JoinPath(a.public, distDir, item.File)
		if err != nil {
			continue
		}
		a.scripts = append(a.scripts, item.File)
		for _, css := range item.Css {
			var err error
			css, err = url.JoinPath(a.public, distDir, css)
			if err != nil {
				continue
			}
			a.styles = append(a.styles, css)
		}
	}
	return nil
}

func (a *assets) mustRead() {
	err := a.read()
	if err != nil {
		panic(err)
	}
}
