package cp

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
)

func TestCookieGet(t *testing.T) {
	tests := []struct {
		name       string
		cookieName string
		cookieVal  string
		getName    string
		expected   string
	}{
		{
			name:       "Getting the value of an existing cookie should return the correct value",
			cookieName: "testCookie",
			cookieVal:  "testVal",
			getName:    "testCookie",
			expected:   "testVal",
		},
		{
			name:       "Getting the value of a non-existent cookie should return an empty string",
			cookieName: "testCookie",
			cookieVal:  "testVal",
			getName:    "nonExistingCookie",
			expected:   "",
		},
	}
	
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: tt.cookieName, Value: tt.cookieVal})
				c := &cookie{
					control: &control{
						request:  req,
						response: httptest.NewRecorder(),
					},
				}
				assert.Equal(t, tt.expected, c.Get(tt.getName))
			},
		)
	}
}

func TestCookieSet(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		expiration time.Duration
	}{
		{
			name:       "TestSetCookieString",
			value:      "Test",
			expiration: 1 * time.Hour,
		},
		{
			name:       "TestSetCookieInt",
			value:      123,
			expiration: 1 * time.Hour,
		},
		{
			name:       "TestSetCookieImmediateExpire",
			value:      "Test",
			expiration: 0 * time.Hour,
		},
	}
	
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				c := &control{
					response: httptest.NewRecorder(),
					request:  httptest.NewRequest(http.MethodGet, "/", nil),
				}
				c.Cookie().Set(test.name, test.value, test.expiration)
				resCookie := c.Response().Raw().Header().Get("Set-Cookie")
				testCookie := resCookie[:strings.Index(resCookie, ";")]
				assert.Equal(t, test.name, testCookie[:strings.Index(testCookie, "=")])
				switch test.value.(type) {
				case int:
					v := testCookie[strings.Index(testCookie, "=")+1:]
					intV, err := strconv.Atoi(v)
					assert.NoError(t, err)
					assert.Equal(t, test.value, intV)
				default:
					assert.Equal(t, test.value, testCookie[strings.Index(testCookie, "=")+1:])
				}
			},
		)
	}
}
