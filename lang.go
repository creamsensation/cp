package cp

import (
	"strings"
	"time"

	"github.com/creamsensation/config"
	"github.com/creamsensation/cookie"
)

type Lang interface {
	Main() string
	Exists() bool
	Current() string
}

type lang struct {
	config  config.Config
	req     Request
	cookie  cookie.Cookie
	current string
}

const (
	langCookieKey = "X-Lang"
)

var (
	langCookieDuration = 365 * 24 * time.Hour
)

func createLang(config config.Config, req Request, cookie cookie.Cookie) *lang {
	l := &lang{
		config: config,
		req:    req,
		cookie: cookie,
	}
	if config.Localization.Enabled && !config.Localization.Path {
		l.current = cookie.Get(langCookieKey)
	}
	if config.Localization.Enabled && config.Localization.Path {
		l.current = l.parseLangFromUrl()
	}
	return l
}

func (l *lang) Main() string {
	var ml config.Language
	for _, item := range l.config.Localization.Languages {
		if !item.Main {
			continue
		}
		ml = item
		break
	}
	if len(ml.Code) == 0 {
		ml = l.config.Localization.Languages[0]
	}
	return ml.Code
}

func (l *lang) Exists() bool {
	return len(l.current) > 0
}

func (l *lang) Current() string {
	return l.current
}

func (l *lang) parseLangFromUrl() string {
	path := l.req.Path()
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
