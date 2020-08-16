package cache

import "net/http"

type Backend interface {
	Get(req *http.Request) *http.Response
	Set(resp *http.Response) *http.Response
	Delete(req *http.Request)
}
