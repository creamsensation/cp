package cp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/tests"
)

func TestState(t *testing.T) {
	type testData struct {
		Value string `json:"value"`
	}
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
	cr.router = createRouter(cr)
	ctrl := createControl(
		cr,
		httptest.NewRequest(http.MethodGet, "/test", nil),
		httptest.NewRecorder(),
	)
	s := createState(ctrl)
	t.Run(
		"set", func(t *testing.T) {
			s.Set(testData{Value: "test"})
			time.Sleep(time.Millisecond * 10)
			assert.True(t, s.Exists())
		},
	)
	t.Run(
		"get", func(t *testing.T) {
			var data testData
			s.Get(&data)
			assert.Equal(t, "test", data.Value)
		},
	)
	t.Run(
		"reset", func(t *testing.T) {
			s.Reset()
			time.Sleep(time.Millisecond * 10)
			assert.False(t, s.Exists())
		},
	)
}
