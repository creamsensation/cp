package cp

import (
	"net/http"
	"sync"
	"time"
	
	"golang.org/x/time/rate"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/form"
	"github.com/creamsensation/gox"
)

func createCsrfMiddleware() func(c Control) Result {
	return func(c Control) Result {
		if c.Request().Is().Get() {
			return c.Continue()
		}
		csrfToken := c.Request().Form().Value(form.CsrfToken)
		csrfName := c.Request().Form().Value(form.CsrfName)
		if len(csrfToken) == 0 {
			return c.Response().Status(http.StatusForbidden).Render(gox.Text("Invalid CSRF token")) // TODO: lepsi errory
		}
		csrf := c.Csrf().Get(csrfToken, csrfName)
		if !csrf.Exist {
			return c.Response().Status(http.StatusForbidden).Render(gox.Text("Invalid CSRF token"))
		}
		if csrf.Name != csrfName || csrf.UserAgent != c.Request().UserAgent() || csrf.Ip != c.Request().Ip() {
			return c.Response().Status(http.StatusForbidden).Render(gox.Text("Invalid CSRF token"))
		}
		c.Csrf().Destroy(csrfToken)
		return c.Continue()
	}
}

func createRateLimitMiddleware(security config.Security) func(c Control) Result {
	type client struct {
		limiter     *rate.Limiter
		lastAttempt time.Time
	}
	
	var mu sync.Mutex
	clients := make(map[string]*client)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastAttempt) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return func(c Control) Result {
		ip := c.Request().Ip()
		if len(ip) == 0 {
			ip = "localhost"
		}
		mu.Lock()
		v, ok := clients[ip]
		if !ok {
			clients[ip] = &client{
				limiter: rate.NewLimiter(
					rate.Every(security.RateLimit.Interval/time.Duration(int64(security.RateLimit.Attempts))),
					security.RateLimit.Attempts,
				),
			}
			v = clients[ip]
		}
		clients[ip].lastAttempt = time.Now()
		if !v.limiter.Allow() {
			mu.Unlock()
			return c.Response().Status(http.StatusTooManyRequests).Text(http.StatusText(http.StatusTooManyRequests))
		}
		mu.Unlock()
		return c.Continue()
	}
}
