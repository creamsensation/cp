package dev

import (
	"time"
	
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/route"
	"github.com/creamsensation/cp/internal/session"
	"github.com/creamsensation/devtool"
	
	"github.com/creamsensation/gox"
)

type Dev interface {
	Log(value ...any) Dev
	Tool() gox.Node
}

type Internal interface {
	Log(value ...any) Dev
	Query(q string, t time.Duration) Dev
	Tool() gox.Node
	SetRoute(route route.Route) Dev
}

type dev struct {
	devtool        *devtool.Devtool
	values         []any
	queries        []devtool.Query
	renderStart    time.Time
	renderEnd      time.Time
	renderDuration time.Duration
	route          route.Route
	session        session.Session
}

func New(tool *devtool.Devtool, session session.Session) Dev {
	return &dev{
		devtool:     tool,
		values:      make([]any, 0),
		queries:     make([]devtool.Query, 0),
		renderStart: time.Now(),
		session:     session,
	}
}

func (d *dev) Log(value ...any) Dev {
	if !env.Development() {
		return d
	}
	d.values = append(d.values, value...)
	return d
}

func (d *dev) Query(q string, t time.Duration) Dev {
	d.queries = append(d.queries, devtool.Query{Value: q, Duration: t})
	return d
}

func (d *dev) SetRoute(route route.Route) Dev {
	d.route = route
	return d
}

func (d *dev) Tool() gox.Node {
	if !env.Development() {
		return gox.Fragment()
	}
	if env.Development() {
		d.calculateRenderTime(time.Now())
	}
	return d.devtool.CreateView(d.renderDuration.String(), d.values, d.queries, d.session.Map(), d.route.Map())
}

func (d *dev) calculateRenderTime(t time.Time) string {
	if d.renderStart.IsZero() {
		return "--"
	}
	d.renderEnd = t
	d.renderDuration = d.renderEnd.Sub(d.renderStart)
	return d.renderDuration.String()
}
