package cp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/cp/internal/tests"
	"github.com/creamsensation/quirk"
)

func TestAuth(t *testing.T) {
	var sessionCookie, tfaCookie *http.Cookie
	db, redis := tests.CreateDatabaseConnection(t), tests.CreateRedisConnection(t)
	DropUsersTable(quirk.New(db))
	CreateUsersTable(quirk.New(db))
	user := User{
		Active:   true,
		Roles:    []string{"owner"},
		Email:    "dominik@linduska.dev",
		Password: "123456789",
	}
	createTestControl := func(request *http.Request, response http.ResponseWriter) *control {
		return &control{
			context: context.Background(),
			config: config.Config{
				Cache: config.Cache{Adapter: cacheAdapter.Redis},
			},
			core: &core{
				databases: map[string]*quirk.DB{
					naming.Main: db,
				},
				redis: redis,
			},
			request:  request,
			response: response,
		}
	}
	t.Run(
		"create user", func(t *testing.T) {
			id := CreateUserManager(db, 0, "").Create(user)
			user.Id = id
			assert.Equal(t, true, id != 0)
		},
	)
	t.Run(
		"auth in", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()
			c := createTestControl(req, res)
			auth := c.Auth().In("dominik@linduska.dev", "123456789")
			cookies := res.Result().Cookies()
			sessionCookie = cookies[0]
			assert.Equal(t, 1, len(cookies))
			assert.Equal(t, true, len(cookies[0].Value) > 0)
			assert.Equal(t, true, auth.Ok())
			assert.Nil(t, auth.Error())
		},
	)
	t.Run(
		"enable tfa", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/enable/tfa", nil)
			res := httptest.NewRecorder()
			req.AddCookie(sessionCookie)
			c := createTestControl(req, res)
			c.Auth().Tfa().Enable()
			u := c.Auth().User().Get()
			assert.Equal(t, true, u.Tfa)
			assert.Equal(t, true, len(u.TfaSecret.String) > 0)
			assert.Equal(t, true, len(u.TfaUrl.String) > 0)
			assert.Equal(t, true, len(u.TfaCodes.String) > 0)
		},
	)
	t.Run(
		"auth in tfa", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/auth/in/tfa", nil)
			res := httptest.NewRecorder()
			c := createTestControl(req, res)
			a := c.Auth().In("dominik@linduska.dev", "123456789")
			cookies := res.Result().Cookies()
			tfaCookie = cookies[0]
			req.AddCookie(tfaCookie)
			assert.Equal(t, true, a.Tfa())
			assert.Equal(t, true, len(tfaCookie.Value) > 0)
			var u User
			c.Cache().Get(tfaCookie.Value, &u)
			assert.Equal(t, true, u.Id > 0)
			c.DB().Q(`SELECT tfa_secret FROM users WHERE id = @id`, quirk.Map{"id": u.Id}).MustExec(&u)
			otp, err := totp.GenerateCode(u.TfaSecret.String, time.Now())
			assert.Nil(t, err)
			_, ok := c.Auth().Tfa().Verify(otp)
			assert.Equal(t, true, ok)
		},
	)
	t.Run(
		"auth out", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/out", nil)
			res := httptest.NewRecorder()
			req.AddCookie(sessionCookie)
			c := createTestControl(req, res)
			session := c.Auth().Session().Get()
			assert.Equal(t, true, session.Id > 0)
			c.Auth().Out()
			time.Sleep(time.Millisecond * 10)
			session = c.Auth().Session().Get()
			assert.Equal(t, true, session.Id < 1)
			sessionCookie = nil
		},
	)
	t.Run(
		"update user", func(t *testing.T) {
			u := User{Active: false}
			CreateUserManager(db, user.Id, "dominik@linduska.dev").Update(u, "active")
			u = CreateUserManager(db, user.Id, "dominik@linduska.dev").Get()
			assert.Equal(t, false, u.Active)
		},
	)
	t.Cleanup(
		func() {
			DropUsersTable(quirk.New(db))
		},
	)
}
