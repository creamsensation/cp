package requester

import (
	"net/http"
	"regexp"

	"github.com/creamsensation/cp/internal/constant/queryKey"
	"github.com/creamsensation/hx"
)

type Is interface {
	Get() bool
	Post() bool
	Put() bool
	Patch() bool
	Delete() bool
	Hx() bool
	Action() bool
	Localized() bool
}

type is struct {
	*http.Request
	localizedUrlMatcher *regexp.Regexp
}

func CreateIs(r *http.Request, localizedUrlMatcher *regexp.Regexp) Is {
	return is{
		Request:             r,
		localizedUrlMatcher: localizedUrlMatcher,
	}
}

func (i is) Get() bool {
	return i.Method == http.MethodGet
}

func (i is) Post() bool {
	return i.Method == http.MethodPost
}

func (i is) Put() bool {
	return i.Method == http.MethodPut
}

func (i is) Patch() bool {
	return i.Method == http.MethodPatch
}

func (i is) Delete() bool {
	return i.Method == http.MethodDelete
}

func (i is) Hx() bool {
	return i.Header.Get(hx.RequestHeaderRequest) == "true"
}

func (i is) Action() bool {
	return i.Hx() && i.Request.URL.Query().Get(queryKey.Action) != ""
}

func (i is) Localized() bool {
	if i.localizedUrlMatcher == nil {
		return false
	}
	return i.localizedUrlMatcher.MatchString(i.URL.Path)
}
