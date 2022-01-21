package httpx

import (
	"net/http"
	"net/url"
	"time"
)

var (
	clientTimeout = 10000 * time.Millisecond
)

func NewClient() *http.Client {
	transport := http.Transport{}
	hc := &http.Client{
		Transport: &transport,
		Timeout:   clientTimeout,
	}

	return hc
}

// NewProxy ...
func NewProxy(proxyUrl *url.URL) *http.Client {
	transport := http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	hc := &http.Client{
		Transport: &transport,
		Timeout:   clientTimeout,
	}
	return hc
}

// NewClientWithTransPort ...
func NewClientWithTransPort(transport *http.Transport) *http.Client {

	hc := &http.Client{
		Transport: transport,
		Timeout:   clientTimeout,
	}
	return hc
}
