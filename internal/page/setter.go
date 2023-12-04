package page

type Setter interface {
	Title(title string) Setter
	Description(description string) Setter
	Keywords(keywords ...string) Setter
	Meta(name, content string) Setter
}

func (s pageSetter) Title(title string) Setter {
	s.page.title = title
	return s
}

func (s pageSetter) Description(description string) Setter {
	s.page.description = description
	return s
}

func (s pageSetter) Keywords(keywords ...string) Setter {
	s.page.keywords = keywords
	return s
}

func (s pageSetter) Meta(name, content string) Setter {
	s.page.metas = append(
		s.page.metas, [2]string{name, content},
	)
	return s
}
