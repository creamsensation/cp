package cp

type Module interface {
	Controllers() Controllers
	Name() string
}

type module struct {
	core   *core
	module Module
}

func createModule(core *core, m Module) *module {
	return &module{
		core:   core,
		module: m,
	}
}

func (m *module) run() {
	m.registerControllersRoutes()
}

func (m *module) registerControllersRoutes() {
	for _, r := range m.module.Controllers() {
		routes := r.Routes()
		for i := range routes {
			routes[i].Module = m.module.Name()
			routes[i].Controller = r.Name()
		}
		m.core.Routes(routes...)
	}
}
