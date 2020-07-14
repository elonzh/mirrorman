package rewrite

import (
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"

	"github.com/elonzh/mirrorman/pkg/config"
)

type Rewriter struct {
	cfg *config.Rewrite
}

func NewRewriter(cfg *config.Rewrite) *Rewriter {
	// TODO: move compile
	for _, rule := range cfg.Rules {
		rule.CompiledPattern = regexp.MustCompile(rule.Pattern)
	}
	return &Rewriter{cfg: cfg}
}

func (r *Rewriter) Register(proxy *goproxy.ProxyHttpServer) {
	proxy.OnRequest().DoFunc(r.Rewrite)
}

func (r *Rewriter) Rewrite(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	originURL := req.URL.String()
	for _, rule := range r.cfg.Rules {
		if rule.CompiledPattern.MatchString(originURL) {
			url := rule.CompiledPattern.ReplaceAllString(originURL, rule.Replace)
			resp := goproxy.NewResponse(req, "", http.StatusFound, "")
			resp.Header.Set("Location", url)
			ctx.Logf("rewrite %s to  %s", originURL, url)
			return req, resp
		}
	}
	return req, nil
}
