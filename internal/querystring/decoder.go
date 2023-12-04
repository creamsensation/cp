package querystring

import (
	"reflect"
	"slices"
	"strconv"
	"time"
	
	"github.com/iancoleman/strcase"
)

type decoder struct {
	*querystring
}

func (r decoder) process() {
	rv := r.rv
	if rv.Type().Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	isPrefix := len(r.prefix) > 0
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		fn := rv.Type().Field(i).Name
		key := strcase.ToKebab(fn)
		if isPrefix {
			key = r.prefix + prefixDivider + key
		}
		if slices.Contains(r.ignore, fv.Type().String()) || slices.Contains(r.ignore, fn) {
			continue
		}
		values, ok := r.request.URL.Query()[key]
		if !ok || len(values) == 0 {
			continue
		}
		if len(values) > 1 && (fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array) {
			v := r.createSlice(fv.Type(), values)
			if v.IsValid() {
				fv.Set(v)
			}
		}
		if len(values) == 1 {
			v := r.createValue(fv.Type(), values[0])
			if v != nil {
				fv.Set(reflect.ValueOf(v))
			}
		}
	}
}

func (r decoder) createSlice(ft reflect.Type, values []string) reflect.Value {
	result := reflect.New(ft)
	for _, value := range values {
		v := r.createValue(ft.Elem(), value)
		if v == nil {
			continue
		}
		result.Elem().Set(reflect.Append(result.Elem(), reflect.ValueOf(v)))
	}
	return result.Elem()
}

func (r decoder) createValue(ft reflect.Type, qv string) any {
	var result any
	switch ft.Kind() {
	case reflect.Struct:
		switch ft {
		case timeType:
			un, err := strconv.Atoi(qv)
			if err == nil {
				result = time.Unix(0, int64(un))
			}
		}
	case reflect.Bool:
		result = qv == "true"
	case reflect.String:
		result = qv
	case reflect.Float32:
		value, err := strconv.ParseFloat(qv, 32)
		if err != nil {
			result = float32(0)
		}
		result = float32(value)
	case reflect.Float64:
		value, err := strconv.ParseFloat(qv, 64)
		if err != nil {
			result = float64(0)
		}
		result = value
	case reflect.Int:
		value, err := strconv.Atoi(qv)
		if err != nil {
			result = 0
		}
		if err == nil {
			result = value
		}
	}
	return result
}
