package cp

import (
	"unicode"
	
	"github.com/creamsensation/cp/internal/util"
)

func Interface[T any]() string {
	return util.GetInterfaceName[T]()
}

func IsFirstCharUpper(v string) bool {
	return unicode.IsUpper(rune(v[0]))
}
