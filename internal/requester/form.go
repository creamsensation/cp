package requester

import "net/http"

type Form interface {
	Value(key string) string
}

type form struct {
	*http.Request
}

func CreateForm(r *http.Request) Form {
	return form{r}
}

func (f form) Value(key string) string {
	return f.FormValue(key)
}
