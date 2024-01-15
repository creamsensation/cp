package assets

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/gox"
)

func TestAssets(t *testing.T) {
	styles := []string{"/.static/.dist/main.css"}
	scripts := []string{"/.static/.dist/main.js"}
	a := New(".static", styles, []string{}, scripts)
	
	t.Run(
		"get", func(t *testing.T) {
			assert.Equal(t, "/.static/test", a.Get("test"))
		},
	)
	
	t.Run(
		"get styles", func(t *testing.T) {
			assert.Equal(
				t, `<link rel="stylesheet" type="text/css" href="/.static/.dist/main.css" />`, gox.Render(a.GetStyles()),
			)
		},
	)
	
	t.Run(
		"get scripts", func(t *testing.T) {
			assert.Equal(
				t, `<script defer src="/.static/.dist/main.js" type="module"></script>`, gox.Render(a.GetScripts()),
			)
		},
	)
	
	t.Run(
		"add style", func(t *testing.T) {
			a.AddStyle("/.static/.dist/test.css")
			assert.Equal(
				t,
				`<link rel="stylesheet" type="text/css" href="/.static/.dist/main.css" /><link rel="stylesheet" type="text/css" href="/.static/.dist/test.css" />`,
				gox.Render(a.GetStyles()),
			)
		},
	)
	
	t.Run(
		"add script", func(t *testing.T) {
			a.AddScript("/.static/.dist/test.js", true)
			assert.Equal(
				t,
				`<script defer src="/.static/.dist/main.js" type="module"></script><script defer src="/.static/.dist/test.js" type="module"></script>`,
				gox.Render(a.GetScripts()),
			)
		},
	)
}
