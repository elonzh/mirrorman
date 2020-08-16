package httptool

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func NewResponse(r *http.Request, status int) *http.Response {
	resp := &http.Response{
		Request:          r,
		TransferEncoding: r.TransferEncoding,
		Header:           make(http.Header),
		StatusCode:       status,
	}
	resp.StatusCode = status
	resp.Status = http.StatusText(status)
	return resp
}

func NewResponseWriter(r *http.Request, status int) *responseWriter {
	resp := NewResponse(r, status)
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
