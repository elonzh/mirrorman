package disk

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/elonzh/mirrorman/pkg/httptool"
)

type Options struct {
	BasePath string
}

func (o *Options) init() {
	if o.BasePath == "" {
		o.BasePath = ".diskcache"
	}
}

type Cache struct {
	opts *Options
}

func (c *Cache) Get(req *http.Request) *http.Response {
	p := httptool.UrlToFilepath(c.opts.BasePath, req.URL)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Debugf("%s is not exist, skip", p)
			return nil
		}
		logrus.Warnf("error when load cache: %v", err)
		return nil
	}
	logrus.Infof("load cache: %s", p)
	rw := httptool.NewResponseWriter(req, http.StatusOK)
	http.ServeFile(rw, req, p)
	return nil
}

func (c *Cache) Exits(req *http.Request) bool {
	p := httptool.UrlToFilepath(c.opts.BasePath, req.URL)
	info, err := os.Stat(p)
	if err == nil && info != nil {
		return true
	}
	return false
}

func (c *Cache) Set(resp *http.Response) *http.Response {
	// TODO: parse filename from header
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
	p := httptool.UrlToFilepath(c.opts.BasePath, resp.Request.URL)
	// TODO: file integrity check
	// TODO: thread safe read and write
	// TODO: if we have too many requests for a file at the same time, the server may run out of disk space
	tee, err := NewTeeFile(resp.Body, p)
	if err != nil {
		logrus.WithError(err).Warnf("error when create file %s", p)
		return resp
	}
	resp.Body = tee
	// TODO: parse filename from header
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
	// TODO: file integrity check
	// TODO: clean up temp files
	return resp
}

func (c *Cache) Delete(req *http.Request) {
	p := httptool.UrlToFilepath(c.opts.BasePath, req.URL)
	err := os.Remove(p)
	if err != nil {
		logrus.WithError(err).Warnln()
	}
}

func NewCache(opts *Options) *Cache {
	if opts == nil {
		opts = new(Options)
	}
	opts.init()
	return &Cache{opts: opts}
}
