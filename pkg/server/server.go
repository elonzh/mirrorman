package server

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/cache"
	"github.com/elonzh/mirrorman/pkg/config"
	"github.com/elonzh/mirrorman/pkg/rewrite"
)

type Server struct {
	proxy *goproxy.ProxyHttpServer
	cfg   *config.Config
}

func (s *Server) Serve() {
	s.proxy.Logger.Printf("server started, listening at %s", s.cfg.Addr)
	s.proxy.Logger.Printf("server exit: %s", http.ListenAndServe(s.cfg.Addr, s.proxy))
}

func (s *Server) register() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	b := cache.NewFsBackend(s.cfg.Cache)
	b.Register(s.proxy)

	r := rewrite.NewRewriter(s.cfg.Rewrite)
	r.Register(s.proxy)
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		proxy: goproxy.NewProxyHttpServer(),
		cfg:   cfg,
	}
	s.proxy.Verbose = s.cfg.Verbose
	if s.proxy.Verbose {
		s.proxy.NonproxyHandler = http.DefaultServeMux
	}
	s.register()
	return s
}
