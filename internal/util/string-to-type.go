package util

import (
	"strconv"
)

func StringToType[T any](v string) T {
	result := *new(T)
	switch any(result).(type) {
	case string:
		result = any(Escape(v)).(T)
	case bool:
		result = any(v == "true").(T)
	case int:
		res, err := strconv.Atoi(v)
		if err == nil {
			result = any(res).(T)
		}
	case float32:
		res, err := strconv.ParseFloat(v, 32)
		if err == nil {
			result = any(res).(T)
		}
	case float64:
		res, err := strconv.ParseFloat(v, 64)
		if err == nil {
			result = any(res).(T)
		}
	}
	return result
}
