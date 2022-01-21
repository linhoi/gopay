package httpx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/linhoi/gopay/common/netx"
)

type Proxy struct {
	rowClient *http.Client

	proxyClient    *http.Client
	proxyTransport *http.Transport

	LoadBalance    interface{}
	CircuitBreaker interface{}

	Config interface{}
}

// NewProxyClient ...
func NewProxyClient(proxyUrl string) (*Proxy, error) {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}

	err = netx.ValidateUrl(proxyUrl)
	if err != nil {
		fmt.Println("proxy is invalid", err)
		return nil, err
	}

	transport := http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	p := &Proxy{
		rowClient:      NewClient(),
		proxyClient:    NewClientWithTransPort(&transport),
		proxyTransport: &transport,
		LoadBalance:    nil,
		CircuitBreaker: nil,
	}
	p.proxyTransport = &transport

	return p, nil
}

// OnProxyChange ...
func (p *Proxy) OnProxyChange(proxyUrl string) {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		fmt.Println("proxy url is invalid", err)
		return
	}

	err = netx.ValidateUrl(proxyUrl)
	if err != nil {
		fmt.Println("proxy is invalid", err)
		return
	}


	transport := http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	// have to reset proxy client, only reset transport is un useful
	p.proxyClient = NewClientWithTransPort(&transport)
	p.proxyTransport = &transport
	fmt.Println("proxy is change to ", proxyUrl)
}

// Get ...
func (p *Proxy) Get(url string) (resp *http.Response, err error) {
	return p.proxyClient.Get(url)
}

func (p *Proxy) Do(req *http.Request) (resp *http.Response, err error) {
	return p.proxyClient.Do(req)
}

// PostForm issues a POST to the specified URL,
// with data's keys and values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and Client.Do.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
func (p *Proxy) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return p.proxyClient.PostForm(url, data)
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// If the provided body is an io.Closer, it is closed after the
// request.
//
// To set custom headers, use NewRequest and Client.Do.
//
// See the Client.Do method documentation for details on how redirects
// are handled.
func (p *Proxy) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return p.proxyClient.Post(url, contentType, body)
}
