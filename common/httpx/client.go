package httpx

import (
	"context"
	"net"
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

// NewTCPProxy ...
func NewTCPProxy(tcpProxy string) (*http.Client, error) {
	dialer := net.Dialer{
		Timeout: clientTimeout,
	}

	// verify if tcpProxy is correct
	conn, err := dialer.DialContext(context.Background(), "tcp", tcpProxy)
	if err != nil {
		return nil, err
	}
	err = conn.Close()
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Timeout: clientTimeout,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, tcpProxy)
			},
		},
	}, nil
}

// NewClientWithTransPort ...
func NewClientWithTransPort(transport *http.Transport) *http.Client {

	hc := &http.Client{
		Transport: transport,
		Timeout:   clientTimeout,
	}
	return hc
}
