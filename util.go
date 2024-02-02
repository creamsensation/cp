package cp

import (
	"unicode"

	"github.com/creamsensation/cp/internal/util"
)

func Interface[T any]() string {
	return util.GetInterfaceName[T]()
}

func IsFirstCharUpper(v string) bool {
	if len(v) == 0 {
		return false
	}
	return unicode.IsUpper(rune(v[0]))
}

func convertMap(vars map[string]string) Map {
	r := make(Map)
	if len(vars) > 0 {
		for k, v := range vars {
			r[k] = v
		}
	}
	return r
}
