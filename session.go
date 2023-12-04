package cp

import (
	"fmt"
	"slices"
	"time"
	
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/cp/internal/constant/cacheKey"
	"github.com/creamsensation/cp/internal/constant/cookieName"
	"github.com/creamsensation/cp/internal/session"
)

type SessionManager interface {
	Exists() bool
	Get() session.Session
	New(user UserReader)
	Renew()
	Destroy()
}

type sessionManager struct {
	*control
}

const (
	defaultSessionDuration = 24 * time.Hour
)

func (s sessionManager) Exists() bool {
	token := s.Cookie().Get(cookieName.Session)
	return len(token) > 0 && s.Cache().Exists(s.getKey(token))
}

func (s sessionManager) Get() session.Session {
	var r session.Session
	token := s.Cookie().Get(cookieName.Session)
	s.Cache().Get(s.getKey(token), &r)
	return r
}

func (s sessionManager) New(user UserReader) {
	token := uniuri.New()
	duration := s.config.Security.Session.Duration
	if duration.Hours() == 0 {
		duration = defaultSessionDuration
	}
	s.Cookie().Set(cookieName.Session, token, duration)
	s.Cache().Set(s.getKey(token), s.createSession(user), duration)
}

func (s sessionManager) Renew() {
	token := s.Cookie().Get(cookieName.Session)
	if len(token) == 0 {
		return
	}
	ss := getSession(s.control)
	if ss.Ip != s.control.Request().Ip() || ss.UserAgent != s.control.Request().UserAgent() {
		return
	}
	duration := s.config.Security.Session.Duration
	if duration.Hours() == 0 {
		duration = defaultSessionDuration
	}
	s.Cookie().Set(cookieName.Session, token, duration)
	s.Cache().Set(s.getKey(token), ss, duration)
}

func (s sessionManager) Destroy() {
	token := s.Cookie().Get(cookieName.Session)
	s.Cookie().Set(cookieName.Session, "", time.Millisecond)
	s.Cache().Set(s.getKey(token), "", time.Millisecond)
}

func (s sessionManager) getKey(token string) string {
	return createSessionKey(token)
}

func (s sessionManager) createSession(user UserReader) session.Session {
	return session.Session{
		Id:        user.GetId(),
		Email:     user.GetEmail(),
		Ip:        s.Request().Ip(),
		UserAgent: s.Request().UserAgent(),
		Roles:     user.GetRoles(),
		Super:     s.containsSuperRole(user.GetRoles()...),
	}
}

func (s sessionManager) containsSuperRole(roles ...string) bool {
	for name, item := range s.config.Security.Role {
		if slices.Contains(roles, name) && item.Super {
			return true
		}
	}
	return false
}

func getSession(c *control) session.Session {
	r := new(session.Session)
	token := c.Cookie().Get(cookieName.Session)
	c.Cache().Get(createSessionKey(token), r)
	return *r
}

func createSessionKey(token string) string {
	return fmt.Sprintf("%s:%s", cacheKey.Session, token)
}
