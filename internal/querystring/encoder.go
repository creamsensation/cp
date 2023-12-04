package querystring

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
	
	"github.com/iancoleman/strcase"
)

type encoder struct {
	*querystring
}

func (r *encoder) process() string {
	rv := r.rv
	result := make([]string, 0)
	if rv.Type().Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		fn := rv.Type().Field(i).Name
		if slices.Contains(r.ignore, fv.Type().String()) || slices.Contains(r.ignore, fn) {
			continue
		}
		if fv.IsZero() {
			continue
		}
		key := strcase.ToKebab(fn)
		if len(r.prefix) > 0 {
			key = r.prefix + prefixDivider + key
		}
		switch fv.Kind() {
		case reflect.Array, reflect.Slice:
			for j := 0; j < fv.Len(); j++ {
				switch fv.Type() {
				case timeType:
					result = append(result, fmt.Sprintf("%s=%d", key, fv.Interface().(time.Time).UnixNano()))
				default:
					result = append(result, fmt.Sprintf("%s=%v", key, fv.Index(j).Interface()))
				}
			}
		default:
			switch fv.Type() {
			case timeType:
				result = append(result, fmt.Sprintf("%s=%d", key, fv.Interface().(time.Time).UnixNano()))
			default:
				result = append(result, fmt.Sprintf("%s=%v", key, fv.Interface()))
			}
		}
	}
	if len(result) == 0 {
		return ""
	}
	return strings.Join(result, "&")
}
