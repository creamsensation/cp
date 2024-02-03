package cp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"
	
	"github.com/dchest/uniuri"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	
	"github.com/creamsensation/cp/internal/constant/cacheKey"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/gox"
	"github.com/creamsensation/quirk"
)

type TfaManager interface {
	GetPendingUserId() int
	PendingVerification() bool
	Active() bool
	Enable(id ...int)
	Disable(id ...int)
	Verify(otp string) (string, bool)
	CreateQrImage(nodes ...gox.Node) gox.Node
}

type tfaManager struct {
	*control
}

func (m tfaManager) GetPendingUserId() int {
	var u User
	token := m.Cookie().Get(cookieName.Tfa)
	m.Cache().Get(cacheKey.Tfa+":"+token, &u)
	return u.Id
}

func (m tfaManager) PendingVerification() bool {
	return len(m.Cookie().Get(cookieName.Tfa)) > 0
}

func (m tfaManager) Active() bool {
	u := m.Auth().User().Get()
	return u.Tfa && len(u.TfaUrl.String) > 0 && len(u.TfaCodes.String) > 0 && len(u.TfaSecret.String) > 0
}

func (m tfaManager) Verify(otp string) (string, bool) {
	token := m.Cookie().Get(cookieName.Tfa)
	if len(token) == 0 {
		return "", false
	}
	var u User
	m.Cache().Get(cacheKey.Tfa+":"+token, &u)
	if u.Id == 0 {
		return "", false
	}
	m.DB().
		Q(fmt.Sprintf(`SELECT id, email, roles, tfa_secret FROM %s`, usersTable)).
		Q("WHERE id = @id", quirk.Map{"id": u.Id}).
		MustExec(&u)
	if valid := totp.Validate(otp, u.TfaSecret.String); !valid {
		return "", false
	}
	m.Cache().Set(token, "", time.Millisecond)
	m.Cookie().Set(cookieName.Tfa, "", time.Millisecond)
	return m.Auth().Session().New(u), true
}

func (m tfaManager) Enable(id ...int) {
	userId := m.getUserId(id...)
	s := m.Auth().User().Get(userId)
	key, err := totp.Generate(
		totp.GenerateOpts{
			Issuer:      m.Request().Host(),
			AccountName: s.Email,
		},
	)
	m.Error().Check(err)
	codes := uniuri.NewLen(tfaCodesSize)
	m.DB().
		Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(
			"SET tfa = @tfa, tfa_codes = @tfa-codes, tfa_secret = @tfa-secret, tfa_url = @tfa-url", quirk.Map{
				"tfa":        true,
				"tfa-codes":  codes,
				"tfa-secret": key.Secret(),
				"tfa-url":    key.String(),
			},
		).
		Q("WHERE id = @id", quirk.Map{"id": userId}).
		MustExec()
}

func (m tfaManager) Disable(id ...int) {
	m.DB().
		Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q("SET tfa = false, tfa_codes = NULL, tfa_secret = NULL, tfa_url = NULL").
		Q("WHERE id = @id", quirk.Map{"id": m.getUserId(id...)}).
		MustExec()
	m.Cookie().Destroy(cookieName.Tfa)
}

func (m tfaManager) CreateQrImage(nodes ...gox.Node) gox.Node {
	var u User
	s := m.Auth().Session().Get()
	err := m.DB().
		Q(fmt.Sprintf(`SELECT tfa_url FROM %s`, usersTable)).
		Q("WHERE id = @id", quirk.Map{"id": u.Id}).
		Exec(&u)
	if err != nil {
		return gox.Text(err)
	}
	key, err := otp.NewKeyFromURL(u.TfaUrl.String)
	if err != nil {
		return gox.Text(err)
	}
	img, err := key.Image(tfaImgSize, tfaImgSize)
	if err != nil {
		return gox.Text(err)
	}
	var buffer bytes.Buffer
	if err = png.Encode(&buffer, img); err != nil {
		return gox.Text(err)
	}
	return gox.Img(
		gox.Src("data:image/png;base64,"+base64.StdEncoding.EncodeToString(buffer.Bytes())),
		gox.Alt(fmt.Sprintf("tfa-qr-user-%d", s.Id)),
		gox.Style(gox.Text(fmt.Sprintf("width:%[1]dpx;height:%[1]dpx;", tfaImgSize))),
		gox.Fragment(nodes...),
	)
}

func (m tfaManager) getUserId(id ...int) int {
	var userId int
	idn := len(id)
	if idn == 0 {
		userId = m.Auth().Session().Get().Id
	}
	if idn > 0 {
		userId = id[0]
	}
	return userId
}
