package cp

import (
	"fmt"
	"time"
	
	"github.com/dchest/uniuri"
	"github.com/matthewhartstonge/argon2"
	
	"github.com/creamsensation/cp/internal/constant/cacheKey"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/quirk"
)

type Auth interface {
	Session() SessionManager
	Tfa() TfaManager
	User(dbname ...string) UserManager
	CustomUser(id int, email string, dbname ...string) UserManager
	
	In(email, password string) AuthIn
	Out()
}

type AuthIn interface {
	Token() string
	Ok() bool
	Tfa() bool
	Error() error
}

type authIn struct {
	token string
	ok    bool
	tfa   bool
	err   error
}

type auth struct {
	*control
}

const (
	tfaImgSize   = 200
	tfaCodesSize = 160
)

func (a auth) Session() SessionManager {
	return sessionManager{a.control}
}

func (a auth) Tfa() TfaManager {
	return tfaManager{a.control}
}

func (a auth) User(dbname ...string) UserManager {
	dbn := naming.Main
	if len(dbname) > 0 {
		dbn = dbname[0]
	}
	s := a.Session().Get()
	return CreateUserManager(a.core.databases[dbn], s.Id, s.Email)
}

func (a auth) CustomUser(id int, email string, dbname ...string) UserManager {
	dbn := naming.Main
	if len(dbname) > 0 {
		dbn = dbname[0]
	}
	return CreateUserManager(a.core.databases[dbn], id, email)
}

func (a auth) In(email, password string) AuthIn {
	var r User
	err := a.DB().
		Q(fmt.Sprintf(`SELECT id, email, roles, password, tfa FROM %s`, usersTable)).
		Q("WHERE email = @email", quirk.Map{"email": email}).
		Q("AND active = true").
		Exec(&r)
	if err != nil {
		return authIn{
			ok:  false,
			err: err,
			tfa: false,
		}
	}
	if r.Id == 0 {
		return authIn{
			ok:  false,
			err: ErrorInvalidCredentials,
			tfa: false,
		}
	}
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(r.Password))
	if !ok || err != nil {
		return authIn{
			ok:  false,
			err: err,
			tfa: false,
		}
	}
	if r.Tfa {
		token := uniuri.New()
		a.Cache().Set(cacheKey.Tfa+":"+token, User{Id: r.Id}, time.Minute*5)
		a.Cookie().Set(cookieName.Tfa, token, time.Minute*5)
		return authIn{
			token: token,
			ok:    true,
			err:   nil,
			tfa:   true,
		}
	}
	return authIn{
		token: a.Session().New(r),
		ok:    true,
		err:   nil,
		tfa:   false,
	}
}

func (a auth) Out() {
	a.Session().Destroy()
}

func (a authIn) Token() string {
	return a.token
}

func (a authIn) Ok() bool {
	return a.ok
}

func (a authIn) Tfa() bool {
	return a.tfa
}

func (a authIn) Error() error {
	return a.err
}
