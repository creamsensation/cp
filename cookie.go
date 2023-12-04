package cp

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/creamsensation/cp/env"
)

type Cookie interface {
	Get(name string) string
	Set(name string, value any, expiration time.Duration)
	Destroy(name string)
}

type cookie struct {
	control *control
}

func (c cookie) Get(name string) string {
	r, err := c.control.request.Cookie(name)
	if err != nil {
		return ""
	}
	return r.Value
}

func (c cookie) Set(name string, value any, expiration time.Duration) {
	http.SetCookie(
		c.control.response, &http.Cookie{
			Name:    name,
			Value:   fmt.Sprintf("%v", value),
			Path:    "/",
			Expires: time.Now().Add(expiration),
			Secure:  env.Production(),
		},
	)
}

func (c cookie) Destroy(name string) {
	c.Set(name, "", time.Millisecond)
}
