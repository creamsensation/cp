package cp

import (
	"fmt"
	"time"
	
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/cp/internal/constant/cookieName"
)

type Notifier interface {
	Get() []Notification
	Add() AddNotification
}

type AddNotification interface {
	Success(title string, subtitle ...string)
	Warning(title string, subtitle ...string)
	Error(err error, subtitle ...string)
}

type notifier struct {
	*control
	name          string
	token         string
	cacheKey      string
	cookieKey     string
	notifications []Notification
}

type Notification struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type notifications struct {
	Notifications []Notification `json:"notifications"`
}

const (
	NotificationTypeSuccess = "success"
	NotificationTypeWarning = "warning"
	NotificationTypeError   = "error"
)

var (
	NotificationExpiration = time.Hour
)

func (n *notifier) Get() []Notification {
	if n.Request().Is().Get() {
		n.load()
	}
	return n.notifications
}

func (n *notifier) Add() AddNotification {
	return n
}

func (n *notifier) Success(title string, subtitle ...string) {
	n.append(NotificationTypeSuccess, title, subtitle...)
}

func (n *notifier) Warning(title string, subtitle ...string) {
	n.append(NotificationTypeWarning, title, subtitle...)
}

func (n *notifier) Error(err error, subtitle ...string) {
	n.append(NotificationTypeError, err.Error(), subtitle...)
}

func (n *notifier) load() {
	n.token = n.Cookie().Get(cookieName.Notification)
	if len(n.token) == 0 {
		n.token = uniuri.New()
	}
	n.cacheKey = n.createCacheKey()
	n.cookieKey = n.createCookieKey()
	n.notifications = n.get()
	n.destroy()
}

func (n *notifier) store() {
	if len(n.notifications) == 0 {
		return
	}
	n.token = uniuri.New()
	n.cacheKey = n.createCacheKey()
	n.cookieKey = n.createCookieKey()
	n.control.Cache().Set(n.cacheKey, notifications{Notifications: n.notifications}, NotificationExpiration)
	n.control.Cookie().Set(n.cookieKey, n.token, NotificationExpiration)
}

func (n *notifier) append(notificationType, title string, subtitle ...string) {
	var sub string
	if len(subtitle) > 0 {
		sub = subtitle[0]
	}
	n.notifications = append(n.notifications, Notification{notificationType, title, sub})
}

func (n *notifier) get() []Notification {
	var r notifications
	n.Cache().Get(n.cacheKey, &r)
	return r.Notifications
}

func (n *notifier) destroy() {
	n.Cache().Destroy(n.cacheKey)
	n.Cookie().Destroy(n.cookieKey)
}

func (n *notifier) createCookieKey() string {
	if len(n.name) == 0 {
		return cookieName.Notification
	}
	return fmt.Sprintf("%s-%s", cookieName.Notification, n.name)
}

func (n *notifier) createCacheKey() string {
	if len(n.name) == 0 {
		return fmt.Sprintf("notification:%s", n.token)
	}
	return fmt.Sprintf("notification-%s:%s", n.name, n.token)
}
