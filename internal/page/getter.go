package page

import "strings"

type Getter interface {
	Title() string
	Description() string
	Keywords() string
	Metas() [][2]string
}

type pageGetter struct {
	*page
}

func (g pageGetter) Title() string {
	return g.title
}

func (g pageGetter) Description() string {
	return g.description
}

func (g pageGetter) Keywords() string {
	return strings.Join(g.keywords, ", ")
}

func (g pageGetter) Metas() [][2]string {
	return g.metas
}
