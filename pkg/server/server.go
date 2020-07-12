package server

import (
	"net/http"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/cache"
	"github.com/elonzh/mirrorman/pkg/config"
)

type Server struct {
	proxy *goproxy.ProxyHttpServer
	cfg   *config.Config
}

func (s *Server) Serve() {
	s.proxy.Logger.Printf("server started, listening at %s", s.cfg.Addr)
	s.proxy.Logger.Printf("server exit: %s", http.ListenAndServe(s.cfg.Addr, s.proxy))
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		proxy: goproxy.NewProxyHttpServer(),
		cfg:   cfg,
	}
	s.proxy.Verbose = s.cfg.Verbose
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	b := cache.NewFsBackend(s.cfg.Cache)
	b.Register(s.proxy)
	return s
}
