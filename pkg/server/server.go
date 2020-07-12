package server

import (
	"log"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/cache"
)

type Server struct {
	Proxy *goproxy.ProxyHttpServer
	Addr  string
}

func (s *Server) Serve() {
	s.Proxy.Logger.Printf("server started, listening at %s", s.Addr)
	log.Fatalln(http.ListenAndServe(s.Addr, s.Proxy))
}

func (s *Server) Init() {
	s.Proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	// https://github.com/gohugoio/hugo/releases/download/v0.73.0/hugo_0.73.0_Linux-64bit.deb
	pattern := regexp.MustCompile(".*mirrors.huaweicloud.com.*")
	s.Proxy.OnRequest(
		goproxy.ReqHostMatches(pattern),
	).DoFunc(cache.LoadCache)
	s.Proxy.OnResponse(
		goproxy.ReqHostMatches(pattern),
	).DoFunc(cache.SaveCache)
}

func NewServer() *Server {
	proxy := goproxy.NewProxyHttpServer()
	return &Server{
		Proxy: proxy,
		Addr:  "",
	}
}
