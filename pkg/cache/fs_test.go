package cache

import (
	"net/url"
	"testing"
)

func Test_urlToFilepath(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://github.com/elazarl/goproxy/archive/v1.1.zip", "https/github.com/elazarl/goproxy/archive/v1.1.zip"},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Fail()
			}
			if got := urlToFilepath("", u); got != tt.want {
				t.Errorf("urlToFilepath() = %v, want %v", got, tt.want)
			}
		})
	}
}
