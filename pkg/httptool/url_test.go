package httptool

import (
	"net/url"
	"testing"
)

func Test_UrlToFilepath(t *testing.T) {
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
			if got := UrlToFilepath("", u); got != tt.want {
				t.Errorf("UrlToFilepath() = %v, want %v", got, tt.want)
			}
		})
	}
}
