package cp

import (
	"errors"
	"fmt"
	"net/http"
	
	"github.com/creamsensation/cp/env"
)

type ErrorHandler interface {
	Check(err error)
	Throw()
	Message(message string, params ...map[string]any) ErrorHandler
	Internal() ErrorHandler
	BadRequest() ErrorHandler
	Forbidden() ErrorHandler
	Unauthorized() ErrorHandler
	NotFound() ErrorHandler
}

type errorHandler struct {
	control *control
	message string
}

type Error struct {
	Message    string
	Route      string
	StatusCode int
}

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorInvalidDatabase    = errors.New("invalid database")
	ErrorForbidden          = errors.New("insufficient rights")
)

func createErrorHandler(control *control) *errorHandler {
	return &errorHandler{
		control: control,
	}
}

func (m *errorHandler) Check(err error) {
	if err != nil {
		if len(m.message) > 0 {
			err = fmt.Errorf("%s: %w", m.message, err)
		}
		if env.Development() {
			fmt.Println(err)
		}
		panic(err)
	}
}

func (m *errorHandler) Throw() {
	panic(errors.New(m.message))
}

func (m *errorHandler) Message(message string, params ...map[string]any) ErrorHandler {
	m.message = m.control.Translate(message, params...)
	return m
}

func (m *errorHandler) Internal() ErrorHandler {
	*m.control.statusCode = http.StatusInternalServerError
	m.message = m.control.Translate("error.internal")
	return m
}

func (m *errorHandler) BadRequest() ErrorHandler {
	*m.control.statusCode = http.StatusBadRequest
	m.message = m.control.Translate("error.bad.request")
	return m
}

func (m *errorHandler) Forbidden() ErrorHandler {
	*m.control.statusCode = http.StatusForbidden
	m.message = m.control.Translate("error.forbidden")
	return m
}

func (m *errorHandler) Unauthorized() ErrorHandler {
	*m.control.statusCode = http.StatusUnauthorized
	m.message = m.control.Translate("error.unauthorized")
	return m
}

func (m *errorHandler) NotFound() ErrorHandler {
	*m.control.statusCode = http.StatusNotFound
	m.message = m.control.Translate("error.not.found")
	return m
}
