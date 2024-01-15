package hx

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/hx"
)

func TestHxResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test/1", nil)
	res := httptest.NewRecorder()
	hxr := New(req, res)
	
	updateIdValue := "#test"
	redirectValue := "/test/2"
	
	hxr.Update(updateIdValue)
	hxr.Redirect(redirectValue)
	
	hxr.PrepareHeaders()
	
	assert.Equal(t, "outerHTML", res.Header().Get(hx.ResponseHeaderReswap))
	assert.Equal(t, updateIdValue, res.Header().Get(hx.ResponseHeaderRetarget))
	assert.Equal(t, redirectValue, res.Header().Get(hx.ResponseHeaderRedirect))
}
