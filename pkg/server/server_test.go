package server

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_extractUrl(t *testing.T) {
	type args struct {
		u      *url.URL
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			"",
			args{&url.URL{
				Scheme: "http",
				Host:   "127.0.0.1",
				Path:   "/proxy/https://example.com/file.txt",
			}, ""},
			nil,
			true,
		},
		{
			"",
			args{&url.URL{
				Scheme:   "http",
				Host:     "127.0.0.1",
				Path:     "/proxy/http://example.com/img/example.png",
				RawQuery: "width=100&height=100",
			}, "/proxy/"},
			&url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "/img/example.png",
				RawQuery: "width=100&height=100",
			},
			false,
		},
		{
			"",
			args{&url.URL{
				Scheme:   "http",
				Host:     "127.0.0.1",
				Path:     "/proxy/http:/example.com/img/example.png",
				RawQuery: "width=100&height=100",
			}, "/proxy/"},
			&url.URL{
				Scheme:   "http",
				Host:     "example.com",
				Path:     "/img/example.png",
				RawQuery: "width=100&height=100",
			},
			false,
		},
		{
			"",
			args{&url.URL{
				Path: "/proxy/https:/mirrors.example.com/helm/v3.2.4/helm-v3.2.4-linux-amd64.tar.gz",
			}, "/proxy/"},
			&url.URL{
				Scheme: "https",
				Host:   "mirrors.example.com",
				Path:   "/helm/v3.2.4/helm-v3.2.4-linux-amd64.tar.gz",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractUrl(tt.args.u, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
