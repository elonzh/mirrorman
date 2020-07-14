package rewrite

import (
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/config"
)

type Rewriter struct {
}

func NewRewriter(cfg *config.Rewrite) *Rewriter {
	return &Rewriter{}
}

func (r *Rewriter) Register(proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest().DoFunc(r.Rewrite)
}

func (r *Rewriter) Rewrite(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// TODO: WIP
	p := regexp.MustCompile("https://storage.googleapis.com.*/golang/(.*)")
	originURL := req.URL.String()
	if !p.MatchString(originURL) {
		ctx.Logf("%s does not match rewrite rule", originURL)
		return req, nil
	}
	url := p.ReplaceAllString(originURL, "https://dl.google.com/go/$1")
	resp := goproxy.NewResponse(req, "", http.StatusFound, "")
	resp.Header.Set("Location", url)
	ctx.Logf("rewrite %s to  %s", originURL, url)
	return req, resp
}
