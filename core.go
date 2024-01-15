package cp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	
	"github.com/go-redis/redis/v8"
	
	"github.com/creamsensation/cp/internal/cache/memory"
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/connect"
	"github.com/creamsensation/cp/internal/constant/cacheAdapter"
	"github.com/creamsensation/cp/internal/filesystem"
	"github.com/creamsensation/cp/internal/handler"
	"github.com/creamsensation/cp/internal/route"
	"github.com/creamsensation/cp/internal/style"
	"github.com/creamsensation/cp/internal/translator"
	"github.com/creamsensation/devtool"
	"github.com/creamsensation/quirk"
)

type Core interface {
	Container(dependencies ...*Dependency) Core
	Controllers(controllers ...Controller) Core
	Middleware(middlewares ...handler.Fn) Core
	Modules(modules ...Module) Core
	Routes(builders ...*route.Builder) Core
	Serve()
	Ui() Ui
}

type core struct {
	assetsReader *assetsReader
	config       config.Config
	databases    map[string]*quirk.DB
	deps         map[string]*Dependency
	devtool      *devtool.Devtool
	fs           *filesystem.FS
	http         *http.Server
	memory       memory.Client
	redis        *redis.Client
	router       *router
	translator   *translator.Translator
	ui           *ui
}

func New(configDir string) Core {
	cfg := config.Parse(configDir)
	wd, _ := os.Getwd()
	c := &core{
		assetsReader: createAssetsReader(wd, cfg.Assets),
		config:       cfg,
		databases:    make(map[string]*quirk.DB),
		deps:         make(map[string]*Dependency),
		devtool:      devtool.New(),
		translator:   translator.New(configDir),
		ui:           createUi(),
	}
	c.onInit()
	c.router = createRouter(c)
	return c
}

func (c *core) Container(dependencies ...*Dependency) Core {
	for _, dep := range dependencies {
		refType := reflect.TypeOf(dep.provider)
		if refType.Kind() != reflect.Ptr {
			fmt.Printf(
				"%s\n",
				style.RedColor.Render(fmt.Sprintf("Dependency [%s] must be ptr", refType.String())),
			)
			continue
		}
		if len(dep.name) == 0 {
			dep.name = refType.String()
		}
		c.deps[dep.name] = dep
	}
	return c
}

func (c *core) Controllers(controllers ...Controller) Core {
	for _, cl := range controllers {
		createController(c, cl).run()
	}
	return c
}

func (c *core) Middleware(middlewares ...handler.Fn) Core {
	c.router.middlewares = append(c.router.middlewares, middlewares...)
	return c
}

func (c *core) Modules(modules ...Module) Core {
	for _, m := range modules {
		createModule(c, m).run()
	}
	return c
}

func (c *core) Routes(builders ...*route.Builder) Core {
	for _, b := range builders {
		route.Process(b, nil, c.config.Languages, c.config.Router)
	}
	c.router.builders = append(c.router.builders, route.CreateFlatBuilders(builders...)...)
	return c
}

func (c *core) Serve() {
	fmt.Printf(
		"🍰 %s [%s] running on port -> :%s \n",
		style.PinkColor.Render("Creampuff"),
		style.GoldColor.Render(c.config.App.Name),
		style.BlueColor.Render(fmt.Sprintf("%d", c.config.App.Port)),
	)
	c.beforeServe()
	err := c.http.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *core) Ui() Ui {
	return c.ui
}

func (c *core) onInit() {
	if c.languagesExist() {
		go c.translator.Prepare()
	}
	go c.createCacheConnection()
	go c.createDatabasesConnections()
	go c.createFilesystem()
}

func (c *core) beforeServe() {
	c.createServer()
}

func (c *core) createServer() {
	c.http = &http.Server{
		Addr:    fmt.Sprintf(":%d", c.config.App.Port),
		Handler: c.router.createHandler(),
	}
}

func (c *core) createCacheConnection() {
	switch c.config.Cache.Adapter {
	case cacheAdapter.Memory:
		c.memory = memory.New(os.TempDir())
	case cacheAdapter.Redis:
		c.redis = connect.Redis(c.config.Cache)
	}
}

func (c *core) createDatabasesConnections() {
	for name, item := range c.config.Database {
		c.databases[name] = connect.Database(item)
	}
}

func (c *core) createFilesystem() {
	cfg := c.config.Filesystem
	c.fs = &filesystem.FS{
		Driver:      cfg.Driver,
		Dir:         cfg.Dir,
		StorageName: cfg.StorageName,
	}
	if len(cfg.Driver) == 0 {
		return
	}
	switch cfg.Driver {
	case filesystem.Local:
		c.fs.Filesystem = filesystem.CreateLocalFilesystem(cfg.Dir)
	case filesystem.Cloud:
		c.fs.Client = connect.CloudFilesystem(cfg)
	}
}

func (c *core) languagesExist() bool {
	for _, l := range c.config.Languages {
		if l.Enabled {
			return true
		}
	}
	return false
}
