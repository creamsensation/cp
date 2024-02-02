package cp

import (
	"fmt"
	"strings"
	"time"

	"github.com/dchest/uniuri"

	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/constant/header"
)

type Csrf interface {
	Get(token, name string) Token
	Create(key, name, ip, userAgent string) string
	Destroy(token string)
	Clean()
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

func (c *csrf) Create(key, name, ip, userAgent string) string {
	token := uniuri.New()
	c.Cookie().Set(cookieName.Csrf+"-"+key, token, c.expiration())
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

func (c *csrf) Clean() {
	cookies := c.Request().Header().Get(header.Cookie)
	if cookies == "" {
		return
	}
	for _, part := range strings.Split(cookies, ";") {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, cookieName.Csrf) {
			continue
		}
		if !strings.Contains(part, "=") {
			continue
		}
		if len(part) < strings.Index(part, "=")+1 {
			continue
		}
		name := part[:strings.Index(part, "=")]
		token := part[strings.Index(part, "=")+1:]
		c.Destroy(token)
		if c.Request().Route() != strings.TrimPrefix(name, cookieName.Csrf+"-") {
			c.Cookie().Destroy(name)
		}
	}
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
