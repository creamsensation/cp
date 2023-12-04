package page

type Page interface {
	Get() Getter
	Set() Setter
}

type page struct {
	title       string
	description string
	keywords    []string
	metas       [][2]string
}

func New() Page {
	return &page{
		keywords: make([]string, 0),
	}
}

type pageSetter struct {
	*page
}

func (p *page) Get() Getter {
	return pageGetter{p}
}

func (p *page) Set() Setter {
	return pageSetter{p}
}
