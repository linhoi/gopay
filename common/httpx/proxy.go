package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/linhoi/gopay/common/breaker"
	"github.com/linhoi/gopay/common/netx"
	"github.com/pkg/errors"
)

type Proxy struct {
	rowClient *http.Client

	proxyClient    *http.Client
	proxyTransport *http.Transport

	LoadBalance interface{}
	// CircuitBreaker is for proxy client
	CircuitBreaker breaker.ICircuitBreaker

	Config interface{}
}

type ProxyConfig struct {
	ProxyUrl      string
	BreakerConfig *breaker.Config
}

var defaultRecordFunc = func(err error) error {
	fmt.Printf("circuit breaker record a error: %s\n", err)
	return nil
}

// NewProxyClient ...
func NewProxyClient(cfg ProxyConfig) (*Proxy, error) {
	proxy, err := url.Parse(cfg.ProxyUrl)
	if err != nil {
		return nil, err
	}

	err = netx.ValidateUrl(cfg.ProxyUrl)
	if err != nil {
		fmt.Println("proxy is invalid", err)
		return nil, err
	}

	transport := http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	var circuitBreaker breaker.ICircuitBreaker
	if cfg.BreakerConfig != nil {
		circuitBreaker, err = breaker.NewHystrixBreaker(*cfg.BreakerConfig)
		if err != nil {
			return nil, err
		}
	}

	p := &Proxy{
		rowClient:      NewClient(),
		proxyClient:    NewClientWithTransPort(&transport),
		proxyTransport: &transport,
		LoadBalance:    nil,
		CircuitBreaker: circuitBreaker,
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
	ok := p.validProxy()
	if !ok {
		return p.rowClient.Get(url)
	}
	return p.do(func() (resp *http.Response, err error) {
		return p.proxyClient.Get(url)
	})
}

// Do ...
func (p *Proxy) Do(req *http.Request) (resp *http.Response, err error) {
	ok := p.validProxy()
	if !ok {
		return p.rowClient.Do(req)
	}
	return p.do(func() (resp *http.Response, err error) {
		return p.proxyClient.Do(req)
	})
}

// PostForm ...
func (p *Proxy) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	ok := p.validProxy()
	if !ok {
		return p.rowClient.PostForm(url, data)
	}
	return p.do(func() (resp *http.Response, err error) {
		return p.proxyClient.PostForm(url, data)
	})

}

// Post ...
func (p *Proxy) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	ok := p.validProxy()
	if !ok {
		return p.rowClient.Post(url, contentType, body)
	}
	return p.do(func() (resp *http.Response, err error) {
		return p.proxyClient.Post(url, contentType, body)
	})
}

func (p *Proxy) validProxy() bool {
	//todo: add more route
	if p.CircuitBreaker == nil {
		fmt.Println("circuit breaker is null")
		return false
	}
	if p.CircuitBreaker.IsOpen() {
		fmt.Println("circuit breaker is open")
		return false
	}
	return true
}

func isContextError(err error) bool {
	coreError := errors.Cause(err)
	return coreError == context.DeadlineExceeded || coreError == context.Canceled
}

func (p *Proxy) do(run func() (resp *http.Response, err error)) (resp *http.Response, err error) {
	breakerErr := p.CircuitBreaker.Do(
		func() error {
			resp, err = run()
			// only mark context.DeadlineExceeded and Canceled as  error
			if isContextError(err) {
				return err
			}
			return nil
		},
		defaultRecordFunc)

	// do not retry in breaker error
	if breakerErr != nil {
		fmt.Printf("breaker error %s\n", breakerErr)
	}
	if resp == nil || err == nil {
		return nil, breakerErr
	}
	return resp, err
}
