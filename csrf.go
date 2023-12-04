package cp

import (
	"fmt"
	"time"
	
	"github.com/dchest/uniuri"
)

type Csrf interface {
	Get(token, name string) Token
	Create(name, ip, userAgent string) string
	Destroy(token string)
}

type csrf struct {
	*control
}

const (
	defaultCsrfExpiration = time.Hour * 24
)

func createCsrf(c *control) *csrf {
	return &csrf{c}
}

type Token struct {
	Exist     bool   `json:"exist"`
	Name      string `json:"name"`
	Ip        string `json:"ip "`
	UserAgent string `json:"userAgent"`
	Token     string `json:"token"`
}

func (c *csrf) Get(token, name string) Token {
	var result Token
	c.Cache().Get(c.createCsrfKey(token), &result)
	if result.Name != name {
		return Token{}
	}
	return result
}

func (c *csrf) Create(name, ip, userAgent string) string {
	token := uniuri.New()
	c.Cache().Set(
		c.createCsrfKey(token),
		Token{
			Exist:     true,
			Name:      name,
			Ip:        ip,
			UserAgent: userAgent,
			Token:     token,
		},
		c.expiration(),
	)
	return token
}

func (c *csrf) Destroy(token string) {
	c.Cache().Destroy(c.createCsrfKey(token))
}

func (c *csrf) createCsrfKey(token string) string {
	return fmt.Sprintf("csrf:%s", token)
}

func (c *csrf) expiration() time.Duration {
	if c.config.Security.Csrf.Duration.Hours() < 1 {
		return defaultCsrfExpiration
	}
	return c.config.Security.Csrf.Duration
}
