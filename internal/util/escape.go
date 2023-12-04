package util

import (
	"strconv"
	"strings"
)

func Escape(value string) string {
	var err error
	replacer := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	value = replacer.Replace(value)
	if strings.Contains(value, "'") || strings.Contains(value, "\"") || strings.Contains(value, "`") {
		value, err = strconv.Unquote(value)
		if err != nil {
			return value
		}
	}
	return value
}
