package route

type Builder struct {
	Route
	Configs        []Config
	Localized      bool
	LocalizedRoute map[string]*Builder
	Subroutes      []*Builder
}

func CreateBuilder(configs ...Config) *Builder {
	return &Builder{
		Localized:      IsLocalized(configs...),
		LocalizedRoute: make(map[string]*Builder),
		Subroutes:      make([]*Builder, 0),
		Configs:        configs,
	}
}

func (b *Builder) Group(builders ...*Builder) *Builder {
	b.Subroutes = append(b.Subroutes, builders...)
	return b
}
