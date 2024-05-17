package cp

import (
	"errors"
)

var (
	ErrorInvalidDatabase = errors.New("invalid database")
	ErrorInvalidLayout   = errors.New("invalid layout")
)

func defaultErrorHandler(c Ctx) error {
	return c.Response().
		Status(c.Response().Intercept().Status()).
		Error(c.Response().Intercept().Error())
}
