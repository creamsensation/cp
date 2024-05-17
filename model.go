package cp

import "net/http"

type Map map[string]any

const (
	Action = "action"
	Main   = "main"
)

const (
	namePrefixDivider = "_"
)

var (
	httpMethods = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodOptions, http.MethodHead, http.MethodConnect, http.MethodTrace,
	}
)
