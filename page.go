package cp

import "strings"

type Page interface {
	Get() PageGetter
	Set() PageSetter
}

type PageGetter interface {
	Title() string
	Description() string
	Keywords() string
	Metas() [][2]string
}

type PageSetter interface {
	Title(title string) PageSetter
	Description(description string) PageSetter
	Keywords(keywords ...string) PageSetter
	Meta(name, content string) PageSetter
}

type page struct {
	title       string
	description string
	keywords    []string
	metas       [][2]string
}

type pageGetter struct {
	*page
}

type pageSetter struct {
	*page
}

func createPage() *page {
	return &page{
		keywords: make([]string, 0),
	}
}

func (p *page) Get() PageGetter {
	return pageGetter{p}
}

func (p *page) Set() PageSetter {
	return pageSetter{p}
}

func (p pageGetter) Title() string {
	return p.title
}

func (p pageGetter) Description() string {
	return p.description
}

func (p pageGetter) Keywords() string {
	return strings.Join(p.keywords, ", ")
}

func (p pageGetter) Metas() [][2]string {
	return p.metas
}

func (p pageSetter) Title(title string) PageSetter {
	p.page.title = title
	return p
}

func (p pageSetter) Description(description string) PageSetter {
	p.page.description = description
	return p
}

func (p pageSetter) Keywords(keywords ...string) PageSetter {
	p.page.keywords = keywords
	return p
}

func (p pageSetter) Meta(name, content string) PageSetter {
	p.page.metas = append(
		p.page.metas, [2]string{name, content},
	)
	return p
}
