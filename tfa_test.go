package cp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/dchest/uniuri"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/header"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/cp/internal/tests"
	"github.com/creamsensation/quirk"
)

func TestTfa(t *testing.T) {
	db := tests.CreateDatabaseConnection(t)
	if db == nil {
		assert.Fail(t, "connection doesn't exist")
		return
	}
	CreateUsersTable(quirk.New(db))
	testCtrl := createControl(
		&core{
			config: config.Config{
				Cache: config.Cache{
					Adapter: cacheAdapter.Redis,
				},
			},
			redis:     tests.CreateRedisConnection(t),
			databases: map[string]*quirk.DB{naming.Main: db},
		},
		httptest.NewRequest(http.MethodGet, "/test", nil),
		httptest.NewRecorder(),
	)
	um := testCtrl.Auth().User()
	um.Create(
		User{
			Active:   true,
			Roles:    []string{"owner"},
			Email:    "dominik@linduska.dev",
			Password: "123456789",
		},
	)
	r := um.Get()
	assert.True(t, r.Id > 0)
	t.Cleanup(
		func() {
			DropUsersTable(quirk.New(db))
		},
	)
	t.Run(
		"pending verification", func(t *testing.T) {
			ctrl := createControl(
				&core{
					config: config.Config{
						Cache: config.Cache{
							Adapter: cacheAdapter.Redis,
						},
					},
					redis:     tests.CreateRedisConnection(t),
					databases: map[string]*quirk.DB{naming.Main: db},
				},
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Tfa, uniuri.New()))
			assert.True(t, ctrl.Auth().Tfa().PendingVerification())
		},
	)
	t.Run(
		"enable", func(t *testing.T) {
			ctrl := createControl(
				&core{
					config: config.Config{
						Cache: config.Cache{
							Adapter: cacheAdapter.Redis,
						},
					},
					redis:     tests.CreateRedisConnection(t),
					databases: map[string]*quirk.DB{naming.Main: db},
				},
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			auth := ctrl.Auth().In("dominik@linduska.dev", "123456789")
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Session, auth.Token()))
			ctrl.Auth().Tfa().Enable()
			auth = ctrl.Auth().In("dominik@linduska.dev", "123456789")
			assert.True(t, auth.Tfa())
		},
	)
	t.Run(
		"verify", func(t *testing.T) {
			ctrl := createControl(
				&core{
					config: config.Config{
						Cache: config.Cache{
							Adapter: cacheAdapter.Redis,
						},
					},
					redis:     tests.CreateRedisConnection(t),
					databases: map[string]*quirk.DB{naming.Main: db},
				},
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			auth := ctrl.Auth().In("dominik@linduska.dev", "123456789")
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Tfa, auth.Token()))
			var u User
			ctrl.DB().Q(
				`SELECT tfa_secret FROM users WHERE email = @email`, quirk.Map{"email": "dominik@linduska.dev"},
			).MustExec(&u)
			otp, err := totp.GenerateCode(u.TfaSecret.String, time.Now())
			assert.Nil(t, err)
			token, tfaOk := ctrl.Auth().Tfa().Verify(otp)
			assert.True(t, len(token) > 0)
			assert.True(t, tfaOk)
		},
	)
	t.Run(
		"disable", func(t *testing.T) {
			ctrl := createControl(
				&core{
					config: config.Config{
						Cache: config.Cache{
							Adapter: cacheAdapter.Redis,
						},
					},
					redis:     tests.CreateRedisConnection(t),
					databases: map[string]*quirk.DB{naming.Main: db},
				},
				httptest.NewRequest(http.MethodGet, "/test", nil),
				httptest.NewRecorder(),
			)
			auth := ctrl.Auth().In("dominik@linduska.dev", "123456789")
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Tfa, auth.Token()))
			var u User
			ctrl.DB().Q(
				`SELECT tfa_secret FROM users WHERE email = @email`, quirk.Map{"email": "dominik@linduska.dev"},
			).MustExec(&u)
			otp, err := totp.GenerateCode(u.TfaSecret.String, time.Now())
			assert.Nil(t, err)
			token, tfaOk := ctrl.Auth().Tfa().Verify(otp)
			assert.True(t, tfaOk)
			ctrl.request.Header.Add(header.Cookie, fmt.Sprintf("%s=%s", cookieName.Session, token))
			ctrl.Auth().Tfa().Disable()
			auth = ctrl.Auth().In("dominik@linduska.dev", "123456789")
			assert.False(t, auth.Tfa())
		},
	)
}
