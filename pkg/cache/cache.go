package cache

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/elonzh/mirrorman/pkg/cache/disk"
	"github.com/elonzh/mirrorman/pkg/config"
)

type ContextData map[string]string

func BackendFromOptions(t string, o map[string]interface{}) (Backend, error) {
	switch t {
	case "disk":
		opts := new(disk.Options)
		err := mapstructure.Decode(o, opts)
		if err != nil {
			return nil, err
		}
		return disk.NewCache(opts), nil
	}
	return nil, fmt.Errorf("unspport cache backend: %s", t)
}

type Cache struct {
	cfg     *config.Cache
	backend Backend
}

func NewCache(cfg *config.Cache) *Cache {
	backend, err := BackendFromOptions(cfg.Backend, cfg.Options)
	if err != nil {
		logrus.WithError(err).Fatalln("error when init cache backend")
	}

	return &Cache{
		cfg:     cfg,
		backend: backend,
	}
}

func (c *Cache) Register(proxy *goproxy.ProxyHttpServer) {
	for _, rule := range c.cfg.Rules {
		reqConditions := make([]goproxy.ReqCondition, 0)
		respConditions := make([]goproxy.RespCondition, 0)
		for _, cond := range rule.Conditions {
			t := cond["type"]
			switch t {
			case "UrlMatches":
				// TODO: check regex
				c := goproxy.UrlMatches(regexp.MustCompile(cond["regex"]))
				reqConditions = append(reqConditions, c)
				respConditions = append(respConditions, c)
			default:
				logrus.Fatalln("unsupported condition type:", t)
			}
		}
		proxy.OnRequest(reqConditions...).DoFunc(c.CacheGet)
		proxy.OnResponse(respConditions...).DoFunc(c.CacheSet)
		logrus.Infof("registered rule: %s", rule.Name)
	}
}

func (c *Cache) CacheGet(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	resp := c.backend.Get(req)
	if resp != nil {
		ctx.UserData = ContextData{"cached": "true"}
		return req, resp
	}
	return req, nil
}

func (c *Cache) checkCacheable(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code, [%s]%s %d", resp.Request.Method, resp.Request.URL, resp.StatusCode)
	}
	if resp.Body == nil {
		return fmt.Errorf("nil response body, [%s]%s %d", resp.Request.Method, resp.Request.URL, resp.StatusCode)
	}
	return nil
}

func (c *Cache) CacheSet(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if ctx.UserData != nil {
		data := ctx.UserData.(ContextData)
		if data["cached"] == "true" {
			logrus.Debugln("get response from cache, skip read remote")
			return resp
		}
	}
	if err := c.checkCacheable(resp); err != nil {
		logrus.WithError(err).Infoln("response is not cacheable")
		return resp
	}
	resp = c.backend.Set(resp)
	return resp
}
