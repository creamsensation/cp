package cp

import "github.com/creamsensation/gox"

type layoutFactory = func(c Ctx, nodes ...gox.Node) gox.Node

type Layout interface {
	Add(name string, layout layoutFactory) Layout
}

type layout struct {
	factories map[string]layoutFactory
}

func createLayout() *layout {
	return &layout{
		factories: make(map[string]layoutFactory),
	}
}

func (l *layout) Add(name string, layout layoutFactory) Layout {
	l.factories[name] = layout
	return l
}
