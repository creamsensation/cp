package cp

import (
	"context"
	"net/http"
	"strings"
	"time"
	
	"github.com/creamsensation/cp/internal/filesystem"
	"github.com/creamsensation/cp/internal/util"
	"github.com/creamsensation/quirk"
	
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/assets"
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/cp/internal/dev"
	"github.com/creamsensation/cp/internal/page"
	"github.com/creamsensation/cp/internal/responder/hx"
	"github.com/creamsensation/cp/internal/result"
	"github.com/creamsensation/cp/internal/route"
)

type Control interface {
	Assets() assets.Assets
	Auth() Auth
	Cache() Cache
	Continue() Result
	Cookie() Cookie
	Config() config.Config
	Create() Create
	Csrf() Csrf
	DB(name ...string) *quirk.Quirk
	Dev() dev.Dev
	Email() Email
	Error() ErrorHandler
	Event() Event
	File() filesystem.Filesystem
	Generate() Generator
	Link(name string, arg ...Map) string
	Notifications() Notifier
	Handle() Handle
	Page() page.Page
	Request() Request
	Response() Response
	State() StateManager
	Switch() Switcher
	Translate(key string, args ...map[string]any) string
}

type control struct {
	assets     assets.Assets
	component  component
	config     config.Config
	core       *core
	dev        dev.Dev
	context    context.Context
	main       *control
	notifier   *notifier
	page       page.Page
	request    *http.Request
	response   http.ResponseWriter
	hxResponse *hx.HxResponse
	route      route.Route
	result     result.Result
	state      StateManager
	statusCode *int
	vars       map[string]string
}

var (
	controlInterfaceName = util.GetInterfaceName[Control]()
)

func createControl(core *core, request *http.Request, response http.ResponseWriter) *control {
	c := &control{
		assets: assets.New(
			core.config.Assets.PublicDir, core.assetter.Styles, make([]string, 0), core.assetter.Scripts,
		),
		config:     core.config,
		context:    context.Background(),
		core:       core,
		page:       page.New(),
		request:    request,
		response:   response,
		hxResponse: hx.New(request, response),
		statusCode: new(int),
		vars:       make(map[string]string),
	}
	c.notifier = &notifier{control: c}
	c.state = createState(c)
	if env.Development() {
		c.dev = dev.New(core.devtool, getSession(c))
	}
	return c
}

func (c *control) Assets() assets.Assets {
	return c.assets
}

func (c *control) Auth() Auth {
	return auth{c}
}

func (c *control) Cache() Cache {
	return createCache(c)
}

func (c *control) Continue() Result {
	return nil
}

func (c *control) Cookie() Cookie {
	return cookie{c}
}

func (c *control) Config() config.Config {
	return c.config
}

func (c *control) Create() Create {
	return create{control: c, component: c.component}
}

func (c *control) Csrf() Csrf {
	return createCsrf(c)
}

func (c *control) DB(name ...string) *quirk.Quirk {
	n := naming.Main
	if len(name) > 0 {
		n = name[0]
	}
	d, ok := c.core.databases[n]
	if !ok {
		panic(ErrorInvalidDatabase)
	}
	q := quirk.New(d)
	q.Subscribe(
		func(q string, t time.Duration) {
			devInternal := c.dev.(dev.Internal)
			devInternal.Query(q, t)
		},
	)
	return q
}

func (c *control) Dev() dev.Dev {
	return c.dev
}

func (c *control) Email() Email {
	return &email{control: c}
}

func (c *control) Error() ErrorHandler {
	return createErrorHandler(c)
}

func (c *control) Event() Event {
	return c.core.event
}

func (c *control) File() filesystem.Filesystem {
	cfg := c.config.Filesystem
	if cfg.Driver == filesystem.Local {
		return c.core.fs.Filesystem
	}
	return filesystem.CreateCloudFilesystem(
		c.context,
		cfg.Dir,
		cfg.StorageName,
		c.core.fs.Client,
	)
}

func (c *control) Generate() Generator {
	return &generator{control: c}
}

func (c *control) Link(name string, arg ...Map) string {
	if c.component != nil && IsFirstCharUpper(name[strings.LastIndex(name, linkLevelDivider)+1:]) {
		return c.Generate().Link().Action(name, arg...)
	}
	return c.Generate().Link().Name(name, arg...)
}

func (c *control) Main() Control {
	return c.main
}

func (c *control) Notifications() Notifier {
	return c.notifier
}

func (c *control) Handle() Handle {
	return handle{control: c}
}

func (c *control) Page() page.Page {
	return c.page
}

func (c *control) Request() Request {
	return request{c}
}

func (c *control) Response() Response {
	return response{c}
}

func (c *control) State() StateManager {
	return c.state
}

func (c *control) Switch() Switcher {
	return switcher{c}
}

func (c *control) Translate(key string, args ...map[string]any) string {
	return c.core.translator.Translate(c.Request().Lang(), key, args...)
}
