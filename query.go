package cp

import "github.com/creamsensation/cp/internal/util"

func Query[T any](c Control, key string) T {
	v := c.Request().Raw().URL.Query().Get(key)
	return util.StringToType[T](v)
}
