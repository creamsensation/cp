package translator

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
	
	"github.com/creamsensation/cp/internal/style"
)

type Translator struct {
	translates map[string]map[string]string
	dir        string
}

const (
	translatesFolder = "translates"
	yamlSuffix       = ".yaml"
)

func New(dir string) *Translator {
	return &Translator{
		translates: make(map[string]map[string]string),
		dir:        fmt.Sprintf("%s/%s", dir, translatesFolder),
	}
}

func (t *Translator) Prepare() {
	t.walk()
}

func (t *Translator) Translate(langCode, key string, args ...map[string]any) string {
	langTranslates, ok := t.translates[langCode]
	if !ok {
		return fmt.Sprintf("<%s:language not found>", langCode)
	}
	translate, ok := langTranslates[key]
	if !ok {
		return key
	}
	if strings.Contains(translate, "{{") && strings.Contains(translate, "}}") && len(args) > 0 {
		translate = strings.Replace(translate, "{{ ", "{{", -1)
		translate = strings.Replace(translate, " }}", "}}", -1)
		for ak, av := range args[0] {
			translate = strings.Replace(translate, fmt.Sprintf("{{%s}}", ak), fmt.Sprintf("%v", av), -1)
		}
	}
	return translate
}

func (t *Translator) walk() {
	if err := filepath.Walk(
		t.dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), yamlSuffix) {
				return errors.New("translates must be *.yaml files")
			}
			lang := strings.TrimSuffix(info.Name(), yamlSuffix)
			if t.translates[lang] == nil {
				t.translates[lang] = make(map[string]string)
			}
			subpath := strings.TrimPrefix(strings.TrimSuffix(path, info.Name()), t.dir)
			subpath = strings.TrimPrefix(subpath, "/")
			subpath = strings.TrimSuffix(subpath, "/")
			pathPrefix := t.createPathPrefix(subpath)
			if err := t.readYaml(lang, path, pathPrefix); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		fmt.Printf("Translates prepared: %s\n", style.RedColor.Render("FAIL"))
		log.Fatalln(err)
	}
	fmt.Printf("Translates prepared: %s\n", style.EmeraldColor.Render("SUCCESS"))
}

func (t *Translator) createPathPrefix(path string) string {
	if len(path) == 0 {
		return ""
	}
	result := make([]string, 0)
	for _, p := range strings.Split(path, "/") {
		if len(p) == 0 {
			continue
		}
		result = append(result, strcase.ToKebab(p))
	}
	return strings.Join(result, ".")
}

func (t *Translator) readYaml(lang, path, pathPrefix string) error {
	yamlBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	yamlData := make(map[string]any)
	if err := yaml.Unmarshal(yamlBytes, yamlData); err != nil {
		return err
	}
	t.parseYaml(lang, pathPrefix, yamlData)
	return nil
}

func (t *Translator) parseYaml(lang, prefix string, data map[string]any) {
	isPrefix := len(prefix) > 0
	for dataKey, item := range data {
		var key string
		if isPrefix {
			key = fmt.Sprintf("%s.%v", prefix, dataKey)
		}
		if !isPrefix {
			key = fmt.Sprintf("%v", dataKey)
		}
		subdata, ok := item.(map[string]any)
		if !ok {
			t.translates[lang][key] = fmt.Sprintf("%v", item)
		}
		if ok {
			t.parseYaml(lang, key, subdata)
		}
	}
}
