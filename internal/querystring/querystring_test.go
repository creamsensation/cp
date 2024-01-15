package querystring

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestQueryString(t *testing.T) {
	type testStruct struct {
		A string
		B int
		C bool
		D float64
		E []string
	}
	qp := map[string]any{
		"a": "test",
		"b": 1,
		"c": true,
		"d": 1.99,
		"e": []string{"a", "b", "c"},
	}
	data := new(testStruct)
	req := httptest.NewRequest("GET", "/test", nil)
	q := req.URL.Query()
	for k, v := range qp {
		switch value := v.(type) {
		case []string:
			for _, item := range value {
				q.Add(k, item)
			}
		default:
			q.Add(k, fmt.Sprintf("%v", v))
		}
	}
	req.URL.RawQuery = q.Encode()
	t.Run(
		"decode", func(t *testing.T) {
			New(data).Request(req).Decode()
			assert.Equal(t, qp["a"], data.A)
			assert.Equal(t, qp["b"], data.B)
			assert.Equal(t, qp["c"], data.C)
			assert.Equal(t, qp["d"], data.D)
			assert.Equal(t, qp["e"], data.E)
		},
	)
	t.Run(
		"encode", func(t *testing.T) {
			vals, err := url.ParseQuery(New(data).Encode())
			assert.Nil(t, err)
			assert.Equal(t, qp["a"], vals.Get("a"))
			assert.Equal(t, fmt.Sprintf("%v", qp["b"]), vals.Get("b"))
			assert.Equal(t, fmt.Sprintf("%v", qp["c"]), vals.Get("c"))
			assert.Equal(t, fmt.Sprintf("%v", qp["d"]), vals.Get("d"))
			assert.Equal(t, qp["e"], vals["e"])
		},
	)
}
