package cp

type Controllers []Controller

type Controller interface {
	Name() string
	Routes() Routes
}

type controller struct {
	core       *core
	controller Controller
}

func createController(core *core, r Controller) *controller {
	return &controller{
		core:       core,
		controller: r,
	}
}

func (l *controller) run() {
	l.registerRoutes()
}

func (l *controller) registerRoutes() {
	routes := l.controller.Routes()
	for i := range routes {
		routes[i].Controller = l.controller.Name()
	}
	l.core.Routes(routes...)
}
