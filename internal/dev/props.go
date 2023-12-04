package dev

import (
	"devtool"
	"github.com/creamsensation/cp/internal/session"
)

type ViewProps struct {
	queries        []devtool.Query
	renderDuration string
	session        session.Session
	values         []any
}
