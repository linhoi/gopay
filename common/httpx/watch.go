package httpx

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/linhoi/gopay/common/balance"
	"github.com/linhoi/gopay/common/breaker"
	"github.com/linhoi/gopay/common/netx"
)

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

// OnCircuitBreakerChange ...
// todo: implement OnCircuitBreakerChange
func (p *Proxy) OnCircuitBreakerChange(c breaker.Config) {
	return
}

// OnBalanceChange ...
func (p *Proxy) OnBalanceChange(c *balance.Config) {
	if c == nil {
		return
	}

	if reflect.DeepEqual(p.Config.BalanceConfig, c) {
		return
	}

	p.LoadBalance.RemoveAll()
	p.setBalance(c)
}
