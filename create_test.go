package cp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/form"
	"github.com/creamsensation/gox"
)

func TestCreate(t *testing.T) {
	c := &create{
		control: &control{
			core: &core{
				router: createRouter(&core{}),
			},
			request: httptest.NewRequest(http.MethodGet, "/test", nil),
		},
	}
	t.Run(
		"component", func(t *testing.T) {
			assert.Equal(t, "test", gox.Render(c.Component(&testComponent{})))
		},
	)
	t.Run(
		"defer", func(t *testing.T) {
			assert.Equal(
				t,
				`<div hx-get="/test" hx-trigger="load" hx-swap="outerHTML" hx-headers='{"HX-Trigger":"load"}'><div>Pending...</div></div>`,
				gox.Render(c.Defer("/test", gox.Div(gox.Text("Pending...")))),
			)
		},
	)
	t.Run(
		"form builder", func(t *testing.T) {
			type testForm struct {
				Test form.Field[string]
			}
			fb := c.FormBuilder()
			fb.Add("test").Id("test").With(form.Text("abc"), form.Validate.Required())
			f, err := form.Build[testForm](fb)
			assert.Nil(t, err)
			assert.Equal(t, "abc", f.Test.Value)
		},
	)
}
