package httpx

import (
	"net/http"
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
