package cp

import (
	"github.com/creamsensation/cp/internal/util"
)

func Var[T any](c Control, key string) T {
	cc, ok := c.(*control)
	if !ok {
		return *new(T)
	}
	v, ok := cc.vars[key]
	if !ok {
		return *new(T)
	}
	return util.StringToType[T](v)
}
