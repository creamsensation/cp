package requester

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/hx"
)

func TestIs(t *testing.T) {
	t.Run(
		"get", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			s := CreateIs(req, nil)
			assert.True(t, s.Get())
		},
	)
	t.Run(
		"post", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			s := CreateIs(req, nil)
			assert.True(t, s.Post())
		},
	)
	t.Run(
		"put", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/test", nil)
			s := CreateIs(req, nil)
			assert.True(t, s.Put())
		},
	)
	t.Run(
		"patch", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, "/test", nil)
			s := CreateIs(req, nil)
			assert.True(t, s.Patch())
		},
	)
	t.Run(
		"delete", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/test", nil)
			s := CreateIs(req, nil)
			assert.True(t, s.Delete())
		},
	)
	t.Run(
		"hx", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			req.Header.Set(hx.RequestHeaderRequest, "true")
			s := CreateIs(req, nil)
			assert.True(t, s.Hx())
		},
	)
	t.Run(
		"localized", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/cs/test/we", nil)
			rgx, err := regexp.Compile(`^/(cs|en)/\b`)
			assert.Nil(t, err)
			s := CreateIs(req, rgx)
			assert.True(t, s.Localized())
		},
	)
	t.Run(
		"non-localized", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test/we", nil)
			rgx, err := regexp.Compile(`^/(cs|en)/\b`)
			assert.Nil(t, err)
			s := CreateIs(req, rgx)
			assert.False(t, s.Localized())
		},
	)
}
