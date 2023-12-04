package cp

import (
	"github.com/creamsensation/gox"
)

type Ui interface {
	ErrorPage(fn uiErrorFn) Ui
	Layout(name string, fn uiLayoutFn) Ui
	Notification(fn uiNotificationFn) Ui
}

type ui struct {
	errorPage    uiErrorFn
	layouts      map[string]uiLayoutFn
	notification uiNotificationFn
}

type uiErrorFn = func(c Control, err Error) gox.Node

type uiLayoutFn = func(c Control, nodes ...gox.Node) gox.Node

type uiNotificationFn = func(c Control, notification Notification) gox.Node

func createUi() *ui {
	return &ui{layouts: make(map[string]uiLayoutFn)}
}

func (u *ui) ErrorPage(fn uiErrorFn) Ui {
	u.errorPage = fn
	return u
}

func (u *ui) Layout(name string, fn uiLayoutFn) Ui {
	u.layouts[name] = fn
	return u
}

func (u *ui) Notification(fn uiNotificationFn) Ui {
	u.notification = fn
	return u
}

func createErrorPage(c *control, statusCode int, err error) string {
	if c.core.ui.errorPage == nil {
		return ""
	}
	return gox.Render(
		c.core.ui.errorPage(
			c, Error{
				Message:    err.Error(),
				Route:      c.route.Name,
				StatusCode: statusCode,
			},
		),
	)
}
