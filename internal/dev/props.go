package dev

import (
	"github.com/creamsensation/cp/internal/session"
	"github.com/creamsensation/devtool"
)

type ViewProps struct {
	queries        []devtool.Query
	renderDuration string
	session        session.Session
	values         []any
}
