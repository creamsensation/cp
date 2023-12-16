package route

import (
	"fmt"
	"net/http"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	cfg "github.com/creamsensation/cp/internal/config"
)

func TestProcessor(t *testing.T) {
	t.Run(
		"process route", func(t *testing.T) {
			path := "test"
			methods := []string{http.MethodGet, http.MethodPost}
			b := CreateBuilder(
				CreateConfig(
					ConfigPath, pathConfig{
						method: methods,
						path:   "/" + path,
					},
				),
			)
			Process(b, nil, cfg.Languages{}, cfg.Router{})
			assert.Equal(t, true, b.Route.Ok)
			assert.Equal(t, fmt.Sprintf(`^/%s/?$`, path), b.Route.Matcher.String())
			assert.Equal(t, "/"+path, b.Route.Path)
			assert.Equal(t, methods, b.Route.Method)
		},
	)
	t.Run(
		"localize non-localized route", func(t *testing.T) {
			path := "test"
			methods := []string{http.MethodGet, http.MethodPost}
			b := CreateBuilder(
				CreateConfig(
					ConfigPath, pathConfig{
						method: methods,
						path:   "/" + path,
					},
				),
				CreateConfig(ConfigLocalize, true),
			)
			Process(b, nil, testsLangs, cfg.Router{})
			for langCode := range testsLangs {
				_, ok := b.LocalizedRoute[langCode]
				assert.Equal(t, true, ok)
				assert.Equal(t, true, b.LocalizedRoute[langCode].Ok)
				assert.Equal(t, fmt.Sprintf(`^/%s/%s/?$`, langCode, path), b.LocalizedRoute[langCode].Matcher.String())
				assert.Equal(t, fmt.Sprintf("/{lang:%s}/%s", langCode, path), b.LocalizedRoute[langCode].Path)
				assert.Equal(t, methods, b.LocalizedRoute[langCode].Method)
			}
		},
	)
	t.Run(
		"localized route", func(t *testing.T) {
			paths := map[string]any{"cs": "clanek", "en": "article"}
			methods := []string{http.MethodGet, http.MethodPost}
			b := CreateBuilder(
				CreateConfig(
					ConfigPath, pathConfig{
						method: methods,
						path:   paths,
					},
				),
			)
			Process(b, nil, testsLangs, cfg.Router{})
			for langCode := range testsLangs {
				_, ok := b.LocalizedRoute[langCode]
				assert.Equal(t, true, ok)
				assert.Equal(t, true, b.LocalizedRoute[langCode].Ok)
				assert.Equal(
					t, fmt.Sprintf(`^/%s/%s/?$`, langCode, paths[langCode]), b.LocalizedRoute[langCode].Matcher.String(),
				)
				assert.Equal(t, fmt.Sprintf("/{lang:%s}/%s", langCode, paths[langCode]), b.LocalizedRoute[langCode].Path)
				assert.Equal(t, methods, b.LocalizedRoute[langCode].Method)
			}
		},
	)
	t.Run(
		"grouped localized route", func(t *testing.T) {
			paths := map[string]any{"cs": "clanek", "en": "article"}
			pathDetail := "{id:[0-9]+}"
			methods := []string{http.MethodGet, http.MethodPost}
			b := CreateBuilder(
				CreateConfig(
					ConfigPath, pathConfig{
						method: methods,
						path:   paths,
					},
				),
			)
			b.Group(
				CreateBuilder(
					CreateConfig(
						ConfigPath, pathConfig{
							method: methods,
							path:   pathDetail,
						},
					),
					CreateConfig(ConfigLocalize, true),
				),
			)
			Process(b, nil, testsLangs, cfg.Router{})
			for langCode := range testsLangs {
				_, ok := b.LocalizedRoute[langCode]
				assert.Equal(t, true, ok)
				assert.Equal(t, true, b.LocalizedRoute[langCode].Ok)
				assert.Equal(
					t, fmt.Sprintf(`^/%s/%s/?$`, langCode, paths[langCode]), b.LocalizedRoute[langCode].Matcher.String(),
				)
				assert.Equal(t, fmt.Sprintf("/{lang:%s}/%s", langCode, paths[langCode]), b.LocalizedRoute[langCode].Path)
				assert.Equal(t, methods, b.LocalizedRoute[langCode].Method)
			}
			for _, item := range b.Subroutes {
				for langCode := range testsLangs {
					_, ok := b.LocalizedRoute[langCode]
					assert.Equal(t, true, ok)
				}
				for langCode, lr := range item.LocalizedRoute {
					assert.Equal(t, true, lr.Ok)
					assert.Equal(
						t, fmt.Sprintf(`^/%s/%s/[0-9]+/?$`, langCode, paths[langCode]), lr.Matcher.String(),
					)
					assert.Equal(t, fmt.Sprintf("/{lang:%s}/%s/{id:[0-9]+}", langCode, paths[langCode]), lr.Path)
					assert.Equal(t, methods, lr.Method)
				}
			}
		},
	)
}
