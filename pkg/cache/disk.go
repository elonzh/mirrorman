package cache

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/elazarl/goproxy"
	"github.com/spf13/viper"
)

type ContextData map[string]string

func urlToFilepath(url *url.URL) string {
	return filepath.Join(viper.GetString("cacheDir"), url.Scheme, url.Host, url.Path)
}
func SaveCache(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if ctx.UserData != nil {
		data := ctx.UserData.(ContextData)
		if data["cached"] == "true" {
			ctx.Logf("get response from cache, skip read remote")
			return resp
		}
	}
	if ctx.Resp.Body == nil {
		ctx.Logf("nil response body: %s", ctx.Req.URL)
		return resp
	}
	// TODO: reduce memory usage
	resp = ctx.Resp
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Logf("Cannot read response %s", err)
		return resp
	}
	_ = resp.Body.Close()

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	// TODO: parse filename from header
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Disposition
	p := urlToFilepath(ctx.Req.URL)
	err = os.MkdirAll(path.Dir(p), os.ModePerm)
	if err != nil {
		ctx.Warnf("error when mkdir: %s", err)
		return resp
	}
	file, err := os.Create(p)
	if err != nil {
		ctx.Warnf("error when creat file %s: %s", p, err)
		return resp
	}
	defer file.Close()
	_, _ = file.Write(b)
	ctx.Logf("file saved: %s", p)
	return resp
}

func LoadCache(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	p := urlToFilepath(req.URL)
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			ctx.Warnf("%s is not exist, skip", p)
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
