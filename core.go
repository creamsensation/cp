package cp

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/creamsensation/config"
)

type Creampuff interface {
	Router
	ErrorHandler(handler Handler) Creampuff
	Layout() Layout
	Run(address string)
	Mux() *http.ServeMux
}

type core struct {
	*router
	*assets
	config       config.Config
	errorHandler Handler
	layout       *layout
	mux          *http.ServeMux
	routes       []*Route
}

const (
	logo = `    _______  ________  ________  ________  ________  ________  ________  ________  ________
  //       \/        \/        \/        \/        \/        \/    /   \/        \/        \
 //        /         /         /         /         /         /         /       __/       __/
/       --/        _/        _/         /         /       __/         /        _/        _/
\________/\____/___/\________/\___/____/\__/__/__/\______/  \________/\_______/ \_______/`
)

func New(cfg config.Config) Creampuff {
	mux := http.NewServeMux()
	rts := make([]*Route, 0)
	c := &core{
		config:       cfg,
		errorHandler: defaultErrorHandler,
		layout:       createLayout(),
		mux:          mux,
		routes:       rts,
	}
	c.router = &router{
		config: cfg,
		mux:    mux,
		prefix: cfg.Router.Prefix,
		routes: &rts,
	}
	c.assets = &assets{
		dir:    cfg.App.Assets,
		public: cfg.App.Public,
	}
	c.router.core = c
	c.router.createGetWildcardRoute()
	c.onInit()
	return c
}

func (c *core) ErrorHandler(handler Handler) Creampuff {
	c.errorHandler = handler
	return c
}

func (c *core) Layout() Layout {
	return c.layout
}

func (c *core) Run(address string) {
	fmt.Println(logo)
	log.Fatalln(http.ListenAndServe(address, c.mux))
}

func (c *core) Mux() *http.ServeMux {
	return c.mux
}

func (c *core) onInit() {
	c.assets.mustRead()
}
