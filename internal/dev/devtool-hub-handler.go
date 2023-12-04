package dev

import (
	"net/http"
	
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/devtool"
)

func CreateDevtoolHubConnectionHandler(
	d *devtool.Devtool, request *http.Request, response http.ResponseWriter,
) bool {
	if !env.Development() || request.URL.Path != devtool.HubPath {
		return false
	}
	err := d.Hub().Connect(request, response)
	if err != nil {
		http.Error(response, "can't upgrade connection for devtool hub", http.StatusInternalServerError)
	}
	return true
}
