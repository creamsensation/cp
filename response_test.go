package cp

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/constant/contentType"
	"github.com/creamsensation/cp/internal/result"
	"github.com/creamsensation/gox"
)

func TestResponse(t *testing.T) {
	cfg := config.Config{
		Cache:     config.Cache{Adapter: cacheAdapter.Redis},
		Languages: config.Languages{"cs": {Enabled: true, Default: true}, "en": {Enabled: true}},
	}
	cr := &core{
		config: cfg,
		ui:     createUi(),
	}
	cr.router = createRouter(cr)
	cr.Routes(
		Route(
			Get("/test"), Name("test"), Handler(
				func(c Control) Result {
					return c.Response().Render(gox.Text("test"))
				},
			),
		),
	)
	cr.router.onInit()
	cr.router.prepareRoutes()
	t.Run(
		"error", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			err := errors.New("test")
			res := any(ctrl.Response().Status(http.StatusNotFound).Error(err)).(result.Result)
			assert.Equal(t, err.Error(), res.Content)
			assert.Equal(t, http.StatusNotFound, res.StatusCode)
		},
	)
	t.Run(
		"file", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			file := bytes.Repeat([]byte("abc"), 5)
			res := any(ctrl.Response().File("test", file)).(result.Result)
			assert.Equal(t, file, res.Data)
			assert.Equal(t, contentType.OctetStream, res.ContentType)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		},
	)
	t.Run(
		"json", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			data := map[string]string{"test": "abc"}
			res := any(ctrl.Response().Json(data)).(result.Result)
			assert.Equal(t, `{"test":"abc"}`, res.Content)
			assert.Equal(t, contentType.Json, res.ContentType)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		},
	)
	t.Run(
		"redirect", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			
			res := any(ctrl.Response().Redirect("test")).(result.Result)
			assert.Equal(t, "/test", res.Content)
			assert.Equal(t, http.StatusFound, res.StatusCode)
		},
	)
	t.Run(
		"render", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			res := any(ctrl.Response().Render(gox.Div(gox.Text("test")))).(result.Result)
			assert.Equal(t, "<div>test</div>", res.Content)
			assert.Equal(t, contentType.Html, res.ContentType)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		},
	)
	t.Run(
		"status", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			ctrl.Response().Status(http.StatusOK)
			res := any(ctrl.response).(*httptest.ResponseRecorder)
			assert.Equal(t, http.StatusOK, res.Code)
		},
	)
	t.Run(
		"text", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			res := any(ctrl.Response().Text("test")).(result.Result)
			assert.Equal(t, "test", res.Content)
			assert.Equal(t, contentType.Text, res.ContentType)
			assert.Equal(t, http.StatusOK, res.StatusCode)
		},
	)
}
