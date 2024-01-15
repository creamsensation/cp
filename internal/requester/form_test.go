package requester

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/constant/contentType"
	"github.com/creamsensation/cp/internal/constant/header"
)

func TestForm(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/test",
		strings.NewReader("a=form"),
	)
	req.Header.Set(header.ContentType, contentType.Form)
	f := CreateForm(req)
	assert.Equal(t, "form", f.Value("a"))
}
