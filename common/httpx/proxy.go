package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/linhoi/gopay/common/balance"
	"github.com/linhoi/gopay/common/breaker"
	"github.com/linhoi/gopay/common/netx"
)

type Proxy struct {
	rowClient *http.Client

	proxyClient    *http.Client
	proxyTransport *http.Transport

	LoadBalance balance.Weight
	// CircuitBreaker is for proxy client
	CircuitBreaker breaker.ICircuitBreaker

	Config *ProxyConfig
}

type ProxyConfig struct {
	ProxyUrl      string
	BreakerConfig *breaker.Config
	BalanceConfig *balance.Config
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
		circuitBreaker, err = breaker.NewHystrixBreaker(cfg.BreakerConfig)
		if err != nil {
			return nil, err
		}
	}

	p := &Proxy{
		rowClient:      NewClient(),
		proxyClient:    NewClientWithTransPort(&transport),
		proxyTransport: &transport,
		CircuitBreaker: circuitBreaker,
	}
	p.proxyTransport = &transport
	p.setBalance(cfg.BalanceConfig)

	return p, nil
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


	if p.LoadBalance == nil {
		return true
	}

	clientType, ok := p.LoadBalance.Next().(string)
	if !ok {
		return false
	}
	if ClientType(clientType) == ClientTypeProxy {
		return true
	}
	return false
}

func isContextError(err error) bool {
	//coreError := errors.Cause(err)
	//return coreError == context.DeadlineExceeded || coreError == context.Canceled
	return strings.Contains(err.Error(), context.Canceled.Error()) || strings.Contains(err.Error(), context.DeadlineExceeded.Error())
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
	if resp == nil && err == nil && breakerErr != nil {
		return nil, breakerErr
	}
	return resp, err
}

func (p *Proxy) setBalance(c *balance.Config) {
	if c == nil || len(c.Items) == 0 {
		return
	}

	if p.LoadBalance == nil {
		p.LoadBalance = balance.NewWeightedRR(c)
	} else {
		p.LoadBalance.OnChange(c)
	}
}
