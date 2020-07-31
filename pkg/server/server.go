package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/cache"
	"github.com/elonzh/mirrorman/pkg/config"
	"github.com/elonzh/mirrorman/pkg/rewrite"
)

const proxyRoute = "/proxy"

type Server struct {
	proxy *goproxy.ProxyHttpServer
	cfg   *config.Config
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	req := request.Clone(nil)
	realUrlStr := strings.TrimLeft(req.URL.Path, "")
	realUrl, err := url.Parse(realUrlStr)
	if err != nil || !realUrl.IsAbs() {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(fmt.Sprintf("invalid url: %s", req.URL.Path)))
		return
	}
	realUrl.RawQuery = req.URL.RawQuery
	req.URL = realUrl
	s.proxy.ServeHTTP(writer, req)
}

func (s *Server) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		logrus.WithField("signal", sig).Infoln("received signal, stopping server")
		cancel()
	}()
	go func() {
		s.proxy.Logger.Printf("http server started, listening at %s", s.cfg.HttpAddr)
		s.proxy.Logger.Printf("http server exit: %s", http.ListenAndServe(s.cfg.HttpAddr, s))
	}()
	go func() {
		s.proxy.Logger.Printf("proxy server started, listening at %s", s.cfg.ProxyAddr)
		s.proxy.Logger.Printf("proxy server exit: %s", http.ListenAndServe(s.cfg.ProxyAddr, s.proxy))
	}()
	<-ctx.Done()
}

func (s *Server) ServeHttp() {
	s.proxy.Logger.Printf("http server started, listening at %s", s.cfg.HttpAddr)
	s.proxy.Logger.Printf("http server exit: %s", http.ListenAndServe(s.cfg.HttpAddr, s))
}

func (s *Server) ServeProxy() {
	s.proxy.Logger.Printf("proxy server started, listening at %s", s.cfg.ProxyAddr)
	s.proxy.Logger.Printf("proxy server exit: %s", http.ListenAndServe(s.cfg.ProxyAddr, s.proxy))
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
