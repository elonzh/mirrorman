package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/cache"
	"github.com/elonzh/mirrorman/pkg/config"
	"github.com/elonzh/mirrorman/pkg/rewrite"
)

const (
	proxyRoute      = "/proxy/"
	shutdownTimeout = 30
)

// negative lookahead isn't supported in golang
var schemePattern = regexp.MustCompile("(https?:/)([^/])")

type Server struct {
	proxy *goproxy.ProxyHttpServer
	cfg   *config.Config
}

func extractUrl(u *url.URL, prefix string) (*url.URL, error) {
	realUrlStr := schemePattern.ReplaceAllString(strings.TrimLeft(u.Path, prefix), "$1/$2")
	realUrl, err := url.Parse(realUrlStr)
	if err != nil || !realUrl.IsAbs() || realUrl.Host == "" {
		return nil, fmt.Errorf("invalid url: %s", u.Path)
	}
	realUrl.RawQuery = u.RawQuery
	return realUrl, nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	req := request.Clone(request.Context())
	realUrl, err := extractUrl(request.URL, proxyRoute)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	req.URL = realUrl
	req.Host = ""
	req.RequestURI = ""
	// http server mod will get all cookies from all sites if we using browser
	delete(req.Header, "Cookie")

	s.proxy.ServeHTTP(writer, req)
}

func (s *Server) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	mux.Handle("/", s)
	proxyServer := http.Server{
		Addr:    s.cfg.ProxyAddr,
		Handler: s.proxy,
	}
	httpServer := http.Server{
		Addr:    s.cfg.HttpAddr,
		Handler: mux,
	}
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signals
		logrus.WithField("signal", sig).Infoln("received signal, stopping server")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			logrus.WithError(err).Errorln("error when shutdown http server")
		}
		ctx, cancel = context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
		defer cancel()
		err = proxyServer.Shutdown(ctx)
		if err != nil {
			logrus.WithError(err).Errorln("error when shutdown proxy server")
		}
		cancel()
	}()
	go func() {
		logrus.Infof("proxy server started, listening at %s", s.cfg.ProxyAddr)
		logrus.Infof("proxy server exit: %s", proxyServer.ListenAndServe())
		cancel()
	}()
	go func() {
		logrus.Infof("http server started, listening at %s", s.cfg.HttpAddr)
		logrus.Infof("http server exit: %s", httpServer.ListenAndServe())
		cancel()
	}()
	<-ctx.Done()
}

func (s *Server) register() {
	s.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	cache.NewCache(s.cfg.Cache).Register(s.proxy)
	rewrite.NewRewriter(s.cfg.Rewrite).Register(s.proxy)
}

func NewServer(cfg *config.Config) *Server {
	s := &Server{
		proxy: goproxy.NewProxyHttpServer(),
		cfg:   cfg,
	}
	s.proxy.Logger = logrus.StandardLogger()
	s.proxy.Verbose = s.cfg.Verbose
	s.register()
	return s
}
