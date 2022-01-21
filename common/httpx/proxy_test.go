package httpx

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewProxy(t *testing.T) {
	type args struct {
		proxyUrl *url.URL
	}

	proxyUrl, err := url.Parse("http://127.0.0.1:1087")
	if err != nil {
		t.Error("proxy url is invalid")
		return
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test with a proxy",
			args: args{proxyUrl: proxyUrl},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewProxy(tt.args.proxyUrl)
			res, err := client.Get(aGoogleUrl)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != http.StatusOK {
				t.Error("res.StatusCode != http.StatusOK")
			}
			t.Log(res)
		})
	}
}

var aGoogleUrl = `https://www.google.com/webhp?tab=ww`

func TestProxyOnProxyChange(t *testing.T) {
	type args struct {
		proxy string
	}
	tests := []struct {
		name string
		args args
		wantErr  bool
	}{
		{
			name: "a in valid proxy",
			args: args{proxy: "http://localhost:8080"},
			wantErr: false,
		},
		{
			name: "a valid proxy",
			args: args{proxy: "http://localhost:108"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewProxyClient("http://127.0.0.1:1087")
			if err != nil {
				return
			}
			res, err := p.Get(aGoogleUrl)
			if err != nil {
				t.Error(err)
			}
			t.Log(res)

			p.OnProxyChange(tt.args.proxy)
			res, err = p.Get(aGoogleUrl)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
			t.Log(res)
		})
	}
}

