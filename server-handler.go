package cp

import (
	"io"
	"net/http"
	"os"
	"strings"
	
	"github.com/andybalholm/brotli"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/contentEncoding"
	"github.com/creamsensation/cp/internal/constant/header"
	"github.com/creamsensation/cp/internal/dev"
	"github.com/creamsensation/cp/internal/route"
)

type serverHandler struct {
	*core
	config        config.Config
	routes        []route.Route
	routesHandler *routesHandler
	staticHandler *staticHandler
}

type routesHandler struct {
	*serverHandler
}

type staticHandler struct {
	*serverHandler
}

type compressedWriter struct {
	http.ResponseWriter
	io.Writer
}

func (w compressedWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func createServerHandler(core *core, routes []route.Route) *serverHandler {
	h := &serverHandler{
		core:   core,
		config: core.config,
		routes: routes,
	}
	h.routesHandler = &routesHandler{
		serverHandler: h,
	}
	h.staticHandler = &staticHandler{
		serverHandler: h,
	}
	return h
}

func (h *serverHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/"+h.config.Assets.PublicPath) {
		h.staticHandler.ServeHTTP(res, req)
		return
	}
	if shouldHandle := dev.CreateDevtoolHubConnectionHandler(h.devtool, req, res); shouldHandle {
		return
	}
	h.routesHandler.ServeHTTP(res, req)
}

func (h *routesHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	br := brotli.NewWriter(res)
	defer func(br *brotli.Writer) {
		_ = br.Close()
	}(br)
	c := createControl(h.serverHandler.core, req, compressedWriter{res, br})
	l := createLifecycle(c)
	res.Header().Set(header.ContentEncoding, contentEncoding.Brotli)
	l.run()
}

func (h *staticHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	br := brotli.NewWriter(res)
	defer func(br *brotli.Writer) {
		_ = br.Close()
	}(br)
	path := h.core.config.Assets.RootPath + req.URL.Path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	res.Header().Set(header.ContentEncoding, contentEncoding.Brotli)
	http.ServeFile(compressedWriter{res, br}, req, path)
}
