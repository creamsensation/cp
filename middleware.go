package cp

import (
	"errors"
	"net/http"
	"slices"
	
	"github.com/creamsensation/firewall"
	"github.com/creamsensation/form"
)

func createLangMiddleware() Handler {
	return func(c Ctx) error {
		cfg := c.Config().Localization
		if !cfg.Enabled || len(cfg.Languages) == 0 {
			return c.Continue()
		}
		if !c.Lang().Exists() {
			mainLangCode := c.Lang().Main()
			c.Cookie().Set(langCookieKey, mainLangCode, langCookieDuration)
			if cfg.Path {
				return c.Response().Redirect("/" + mainLangCode + "/")
			}
		}
		var langExists bool
		current := c.Lang().Current()
		for _, l := range cfg.Languages {
			if l.Code == current {
				langExists = true
			}
		}
		if !langExists && cfg.Path {
			return c.Response().Redirect("/" + c.Lang().Main() + "/")
		}
		return c.Continue()
	}
}

func createCsrfMiddleware() Handler {
	return func(c Ctx) error {
		if c.Request().Is().Get() {
			if c.Request().Is().Action() {
				return c.Continue()
			}
			if err := c.Csrf().Clean(c.Request().Name()); err != nil {
				return c.Response().Refresh()
			}
			return c.Continue()
		}
		token := c.Request().Form().Get(form.CsrfToken)
		name := c.Request().Form().Get(form.CsrfName)
		if len(token) == 0 {
			return c.Response().Refresh()
		}
		t, err := c.Csrf().Get(name, token)
		if err != nil {
			return err
		}
		if !t.Exists {
			return c.Response().Refresh()
		}
		if t.Name != name || t.UserAgent != c.Request().UserAgent() || t.Ip != c.Request().Ip() {
			return c.Response().Refresh()
		}
		if err := c.Csrf().Destroy(t); err != nil {
			return c.Response().Refresh()
		}
		return c.Continue()
	}
}

func createFirewallMiddleware(firewalls []firewall.Firewall) Handler {
	return func(c Ctx) error {
		if len(firewalls) == 0 {
			return c.Continue()
		}
		var err error
		session, err := c.Auth().Session().Get()
		if err != nil {
			return c.Response().Error(err)
		}
		if session.Super {
			if err := c.Auth().Session().Renew(); err != nil {
				return c.Response().Status(http.StatusInternalServerError).Error(err)
			}
		}
		results := make([]firewall.Result, len(firewalls))
		secret := c.Request().Header().Get("secret")
		if !session.Super {
			for i, f := range firewalls {
				sessionRoles := make([]firewall.Role, 0)
				for _, sr := range session.Roles {
					for _, fr := range f.Roles {
						if slices.ContainsFunc(
							sessionRoles, func(r firewall.Role) bool {
								return r.Name == fr.Name
							},
						) {
							continue
						}
						if fr.Name != sr {
							continue
						}
						sessionRoles = append(sessionRoles, fr)
					}
				}
				results[i] = f.Try(
					firewall.Attempt{
						Secret: secret,
						Roles:  sessionRoles,
					},
				)
			}
		}
		var allowed bool
		var redirect string
		for _, r := range results {
			if r.Ok {
				allowed = true
				continue
			}
			if len(r.Redirect) > 0 {
				redirect = r.Redirect
			}
			if r.Err != nil {
				err = r.Err
			}
		}
		redirectExists := len(redirect) > 0
		errorExists := err != nil
		if !allowed && !redirectExists && !errorExists {
			return c.Response().Status(http.StatusForbidden).Error(errors.New(http.StatusText(http.StatusForbidden)))
		}
		if !allowed && redirectExists {
			return c.Response().Redirect(redirect)
		}
		if !allowed && errorExists {
			return c.Response().Error(err)
		}
		if allowed {
			if err := c.Auth().Session().Renew(); err != nil {
				return c.Response().Status(http.StatusInternalServerError).Error(err)
			}
		}
		return c.Continue()
	}
}
