package hx

import (
	"fmt"
	"net/http"
	"strings"
	
	. "github.com/creamsensation/hx"
)

type Response interface {
	Location(url string) Response
	PushUrl(url string) Response
	Redirect(url string) Response
	Refresh(condition bool) Response
	ReplaceUrl(url string) Response
	Trigger(event ...string) Response
	TriggerAfterSettle(event ...string) Response
	TriggerAfterSwap(event ...string) Response
	Update(id string) Options
	Exists(target ...string) bool
}

type HxResponse struct {
	request            *http.Request
	response           http.ResponseWriter
	options            *options
	location           string
	pushUrl            string
	redirect           string
	refresh            string
	replaceUrl         string
	trigger            string
	triggerAfterSettle string
	triggerAfterSwap   string
}

func New(request *http.Request, response http.ResponseWriter) *HxResponse {
	r := &HxResponse{
		request:  request,
		response: response,
		options:  &options{},
	}
	return r
}

func (r *HxResponse) Location(url string) Response {
	r.location = url
	return r
}

func (r *HxResponse) PushUrl(url string) Response {
	r.pushUrl = url
	return r
}

func (r *HxResponse) Redirect(url string) Response {
	r.redirect = url
	return r
}

func (r *HxResponse) Refresh(condition bool) Response {
	r.refresh = fmt.Sprintf("%t", condition)
	return r
}

func (r *HxResponse) ReplaceUrl(url string) Response {
	r.replaceUrl = url
	return r
}

func (r *HxResponse) Trigger(event ...string) Response {
	r.trigger = strings.Join(event, " ")
	return r
}

func (r *HxResponse) TriggerAfterSettle(event ...string) Response {
	r.triggerAfterSettle = strings.Join(event, " ")
	return r
}

func (r *HxResponse) TriggerAfterSwap(event ...string) Response {
	r.triggerAfterSwap = strings.Join(event, " ")
	return r
}

func (r *HxResponse) Exists(target ...string) bool {
	if len(target) > 0 {
		return r.options.target == "#"+target[0]
	}
	return len(r.options.target) > 0
}

func (r *HxResponse) Options() bool {
	return len(r.options.target) > 0
}

func (r *HxResponse) Update(id string) Options {
	id = strings.TrimSpace(id)
	if !strings.HasPrefix(id, "#") {
		id = "#" + id
	}
	r.options.target = id
	return r.options
}

func (r *HxResponse) PrepareHeaders() {
	r.prepareHistoryHeaders()
	r.prepareHtmlUpdateHeaders()
}

func (r *HxResponse) prepareHistoryHeaders() {
	if len(r.location) > 0 {
		r.response.Header().Set(ResponseHeaderLocation, r.location)
	}
	if len(r.pushUrl) > 0 {
		r.response.Header().Set(ResponseHeaderPushUrl, r.pushUrl)
	}
	if len(r.redirect) > 0 {
		r.response.Header().Set(ResponseHeaderRedirect, r.redirect)
	}
	if len(r.refresh) > 0 {
		r.response.Header().Set(ResponseHeaderRefresh, r.refresh)
	}
	if len(r.replaceUrl) > 0 {
		r.response.Header().Set(ResponseHeaderReplaceUrl, r.replaceUrl)
	}
	if len(r.trigger) > 0 {
		r.response.Header().Set(ResponseHeaderTrigger, r.trigger)
	}
	if len(r.triggerAfterSettle) > 0 {
		r.response.Header().Set(ResponseHeaderTriggerAfterSettle, r.triggerAfterSettle)
	}
	if len(r.triggerAfterSwap) > 0 {
		r.response.Header().Set(ResponseHeaderTriggerAfterSwatp, r.triggerAfterSwap)
	}
}

func (r *HxResponse) prepareHtmlUpdateHeaders() {
	if len(r.options.target) == 0 {
		return
	}
	r.response.Header().Set(ResponseHeaderRetarget, r.options.target)
	swap := make([]string, 0)
	if len(r.options.swap) == 0 {
		swap = append(swap, SwapOuterHtml)
	}
	if len(r.options.swap) > 0 {
		swap = append(swap, r.options.swap)
	}
	if len(r.options.modifier) > 0 {
		swap = append(swap, r.options.modifier)
	}
	r.response.Header().Set(ResponseHeaderReswap, strings.Join(swap, " "))
}
