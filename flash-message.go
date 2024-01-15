package cp

import (
	"fmt"
	"time"
	
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/cp/internal/constant/cookieName"
)

type FlashMessenger interface {
	Get() []FlashMessage
	Add() AddFlashMessage
}

type AddFlashMessage interface {
	Success(title string, subtitle ...string)
	Warning(title string, subtitle ...string)
	Error(err error, subtitle ...string)
}

type flashMessenger struct {
	*control
	name      string
	token     string
	cacheKey  string
	cookieKey string
	messages  []FlashMessage
}

type FlashMessage struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type flashMessengers struct {
	Messages []FlashMessage `json:"messages"`
}

const (
	FlashMessageTypeSuccess = "success"
	FlashMessageTypeWarning = "warning"
	FlashMessageTypeError   = "error"
)

var (
	FlashMessageExpiration = time.Hour
)

func (n *flashMessenger) Get() []FlashMessage {
	if n.Request().Is().Get() {
		n.load()
	}
	return n.messages
}

func (n *flashMessenger) Add() AddFlashMessage {
	return n
}

func (n *flashMessenger) Success(title string, subtitle ...string) {
	n.append(FlashMessageTypeSuccess, title, subtitle...)
}

func (n *flashMessenger) Warning(title string, subtitle ...string) {
	n.append(FlashMessageTypeWarning, title, subtitle...)
}

func (n *flashMessenger) Error(err error, subtitle ...string) {
	n.append(FlashMessageTypeError, err.Error(), subtitle...)
}

func (n *flashMessenger) load() {
	if len(n.token) == 0 {
		n.token = n.Cookie().Get(cookieName.Flash)
	}
	if len(n.token) == 0 {
		n.token = uniuri.New()
	}
	n.cacheKey = n.createCacheKey()
	n.cookieKey = n.createCookieKey()
	n.messages = n.get()
	n.destroy()
}

func (n *flashMessenger) store() {
	if len(n.messages) == 0 {
		return
	}
	n.token = uniuri.New()
	n.cacheKey = n.createCacheKey()
	n.cookieKey = n.createCookieKey()
	n.control.Cache().Set(n.cacheKey, flashMessengers{Messages: n.messages}, FlashMessageExpiration)
	n.control.Cookie().Set(n.cookieKey, n.token, FlashMessageExpiration)
}

func (n *flashMessenger) append(notificationType, title string, subtitle ...string) {
	var sub string
	if len(subtitle) > 0 {
		sub = subtitle[0]
	}
	n.messages = append(n.messages, FlashMessage{notificationType, title, sub})
}

func (n *flashMessenger) get() []FlashMessage {
	var r flashMessengers
	n.Cache().Get(n.cacheKey, &r)
	return r.Messages
}

func (n *flashMessenger) destroy() {
	n.Cache().Destroy(n.cacheKey)
	n.Cookie().Destroy(n.cookieKey)
}

func (n *flashMessenger) createCookieKey() string {
	if len(n.name) == 0 {
		return cookieName.Flash
	}
	return fmt.Sprintf("%s-%s", cookieName.Flash, n.name)
}

func (n *flashMessenger) createCacheKey() string {
	if len(n.name) == 0 {
		return fmt.Sprintf("flash:%s", n.token)
	}
	return fmt.Sprintf("flash-%s:%s", n.name, n.token)
}
