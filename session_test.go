package cp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/header"
	"github.com/creamsensation/cp/internal/tests"
)

func TestSession(t *testing.T) {
	var token string
	cfg := config.Config{
		Cache: config.Cache{Adapter: cacheAdapter.Redis},
		Security: config.Security{
			Role: map[string]config.SecurityRole{
				"owner": {Super: true},
			},
			Session: config.SecuritySession{Duration: time.Hour * 24},
		},
	}
	cr := &core{
		config: cfg,
		redis:  tests.CreateRedisConnection(t),
	}
	t.Run(
		"new", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			sm := &sessionManager{control: ctrl}
			token = sm.New(
				User{
					Id:    1,
					Roles: []string{"owner"},
					Email: "dominik@linduska.dev",
				},
			)
			assert.True(t, len(token) > 0)
		},
	)
	t.Run(
		"exists", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Session, token))
			sm := &sessionManager{control: ctrl}
			assert.True(t, sm.Exists())
		},
	)
	t.Run(
		"get", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Session, token))
			sm := &sessionManager{control: ctrl}
			s := sm.Get()
			assert.Equal(t, 1, s.Id)
			assert.Equal(t, []string{"owner"}, s.Roles)
			assert.Equal(t, "dominik@linduska.dev", s.Email)
		},
	)
	t.Run(
		"renew", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Session, token))
			sm := &sessionManager{control: ctrl}
			sm.Renew()
			time.Sleep(time.Millisecond * 10)
			assert.True(t, sm.Exists())
			assert.True(t, strings.Contains(ctrl.response.Header().Get(header.SetCookie), cookieName.Session))
		},
	)
	t.Run(
		"destroy", func(t *testing.T) {
			ctrl := createControl(
				cr,
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Session, token))
			sm := &sessionManager{control: ctrl}
			sm.Destroy()
			time.Sleep(time.Millisecond * 10)
			assert.False(t, sm.Exists())
			assert.True(t, strings.Contains(ctrl.response.Header().Get(header.SetCookie), cookieName.Session))
		},
	)
}
