package cache

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/elazarl/goproxy"
	"github.com/sirupsen/logrus"

	"github.com/elonzh/mirrorman/pkg/config"
)

type ContextData map[string]string

type FsBackend struct {
	cfg *config.Cache
	log goproxy.Logger
}

func NewFsBackend(cfg *config.Cache) *FsBackend {
	return &FsBackend{
		cfg: cfg,
		// TODO: configurable logger
		log: logrus.New(),
	}
}

func (b *FsBackend) Register(proxy *goproxy.ProxyHttpServer) {
	b.log = proxy.Logger
	for _, rule := range b.cfg.Rules {
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
				b.log.Printf("unspported condition type: %s", t)
				os.Exit(1)
			}
		}
		proxy.OnRequest(reqConditions...).DoFunc(b.CacheGet)
		proxy.OnResponse(respConditions...).DoFunc(b.CacheSet)
		b.log.Printf("registered rule: %s", rule.Name)
	}
}

func (b *FsBackend) CacheGet(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	p := urlToFilepath(b.cfg.Dir, req.URL)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.Logf("%s is not exist, skip", p)
			return req, nil
		}
		ctx.Warnf("error when load cache: %v", err)
		return req, nil
	}
	ctx.Logf("load cache: %s", p)
	rw := NewResponseWriter(req, http.StatusOK)
	http.ServeFile(rw, req, p)
	ctx.UserData = ContextData{"cached": "true"}
	return req, rw.Response
}

func (b *FsBackend) checkCacheable(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("invalid status code, [%s]%s %d", resp.Request.Method, resp.Request.URL, resp.StatusCode)
	}
	if resp.Body == nil {
		return fmt.Errorf("nil response body, [%s]%s %d", resp.Request.Method, resp.Request.URL, resp.StatusCode)
	}
	return nil
}

func (b *FsBackend) CacheSet(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if ctx.UserData != nil {
		data := ctx.UserData.(ContextData)
		if data["cached"] == "true" {
			ctx.Logf("get response from cache, skip read remote")
			return resp
		}
	}
	if err := b.checkCacheable(resp); err != nil {
		ctx.Warnf("response is not cacheable: %s", err)
		return resp
	}
	resp = ctx.Resp
	// TODO: parse filename from header
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
	p := urlToFilepath(b.cfg.Dir, ctx.Req.URL)
	// TODO: file integrity check
	// TODO: thread safe read and write
	// TODO: if we have too many requests for a file at the same time, the server may run out of disk space
	tee, err := newTeeFile(resp.Body, p)
	if err != nil {
		ctx.Warnf("error when creat file %s: %s", p, err)
		return resp
	}
	resp.Body = tee
	// TODO: parse filename from header
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
	// TODO: file integrity check
	// TODO: clean up temp files
	return resp
}

func urlToFilepath(baseDir string, url *url.URL) string {
	return filepath.Join(baseDir, url.Scheme, url.Host, url.Path)
}

func newTeeFile(r io.ReadCloser, filename string) (io.ReadCloser, error) {
	file, err := ioutil.TempFile("", "mirrorman_*.download")
	if err != nil {
		return nil, err
	}
	return &teeFile{
		r:       r,
		tmpFile: file,
		dst:     filename,
	}, nil
}

type teeFile struct {
	r       io.ReadCloser
	tmpFile *os.File
	dst     string
}

func (t *teeFile) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.tmpFile.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *teeFile) Close() error {
	err := t.r.Close()
	if err != nil {
		return err
	}
	log.Println("close reader")
	err = t.tmpFile.Close()
	log.Println("close tmpFile:", t.tmpFile.Name(), err)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(t.dst), os.ModePerm)
	log.Println("MkdirAll:", filepath.Dir(t.dst), err)
	if err != nil {
		return err
	}
	err = os.Rename(t.tmpFile.Name(), t.dst)
	log.Println("rename tmpFile:", t.tmpFile.Name(), t.dst, err)
	if err != nil {
		return err
	}
	return nil
}

func NewResponseWriter(r *http.Request, status int) *responseWriter {
	resp := &http.Response{
		Request:          r,
		TransferEncoding: r.TransferEncoding,
		Header:           make(http.Header),
		StatusCode:       status,
	}
	resp.StatusCode = status
	resp.Status = http.StatusText(status)
	buf := new(bytes.Buffer)
	resp.Body = ioutil.NopCloser(buf)
	return &responseWriter{
		Response: resp,
		buf:      buf,
	}
}

type responseWriter struct {
	Response *http.Response
	buf      *bytes.Buffer
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.buf.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.Response.StatusCode = statusCode
}

func (r *responseWriter) Header() http.Header {
	return r.Response.Header
}
