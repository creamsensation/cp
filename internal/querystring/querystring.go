package querystring

import (
	"net/http"
	"reflect"
)

type QueryString interface {
	IgnoreInterface(names ...string) QueryString
	Request(req *http.Request) QueryString
	Prefix(prefix string) QueryString
	Decode()
	Encode() string
}

type querystring struct {
	value   any
	rv      reflect.Value
	request *http.Request
	ignore  []string
	prefix  string
}

const (
	prefixDivider = "_"
)

func New(value any) QueryString {
	return &querystring{
		value: value,
		rv:    reflect.ValueOf(value),
	}
}

func (q *querystring) IgnoreInterface(names ...string) QueryString {
	q.ignore = append(q.ignore, names...)
	return q
}

func (q *querystring) Prefix(prefix string) QueryString {
	q.prefix = prefix
	return q
}

func (q *querystring) Request(req *http.Request) QueryString {
	q.request = req
	return q
}

func (q *querystring) Decode() {
	decoder{q}.process()
}

func (q *querystring) Encode() string {
	e := &encoder{q}
	return e.process()
}
