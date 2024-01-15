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
	
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/gox"
)

type TfaManager interface {
	PendingVerification() bool
	Enable()
	Disable()
	Verify(otp string) (string, bool)
	CreateQrImage() gox.Node
}

type tfaManager struct {
	*control
}

func (m tfaManager) PendingVerification() bool {
	return len(m.Cookie().Get(cookieName.Tfa)) > 0
}

func (m tfaManager) Enable() {
	s := m.Auth().Session().Get()
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
		Q("SET tfa = ?, tfa_codes = ?, tfa_secret = ?, tfa_url = ?", true, codes, key.Secret(), key.String()).
		Q("WHERE id = ?", s.Id).
		MustExec()
}

func (m tfaManager) Verify(otp string) (string, bool) {
	token := m.Cookie().Get(cookieName.Tfa)
	if len(token) == 0 {
		return "", false
	}
	var u User
	m.Cache().Get(token, &u)
	if u.Id == 0 {
		return "", false
	}
	m.DB().
		Q(fmt.Sprintf(`SELECT id, email, roles, tfa_secret FROM %s`, usersTable)).
		Q("WHERE id = ?", u.Id).
		MustExec(&u)
	if valid := totp.Validate(otp, u.TfaSecret.String); !valid {
		return "", false
	}
	m.Cache().Set(token, "", time.Millisecond)
	m.Cookie().Set(cookieName.Tfa, "", time.Millisecond)
	return m.Auth().Session().New(u), true
}

func (m tfaManager) Disable() {
	s := m.Auth().Session().Get()
	m.DB().
		Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q("SET tfa = ?, tfa_codes = NULL, tfa_secret = NULL, tfa_url = NULL", false).
		Q("WHERE id = ?", s.Id).
		MustExec()
}

func (m tfaManager) CreateQrImage() gox.Node {
	var u User
	s := m.Auth().Session().Get()
	err := m.DB().
		Q(fmt.Sprintf(`SELECT tfa_url FROM %s`, usersTable)).
		Q("WHERE id = ?", s.Id).
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
	)
}
