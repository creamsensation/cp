package cp

import "errors"

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
