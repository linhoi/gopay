package netx

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	clientTimeout = 10000 * time.Millisecond
)

// ValidateUrl ...
/*
If you want to check if a Web server answers on a certain URL,
you can invoke an HTTP GET request using net/http.
You will get a timeout if the server doesn't response at all.
You might also check the response status.
*/
func ValidateUrl(urlStr string) error {
	_, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c := http.Client{
		Timeout: clientTimeout,
	}
	_, err = c.Get(urlStr)
	if err != nil {
		return err
	}
	c.CloseIdleConnections()

	return nil

}

// ValidateAddress ...
func ValidateAddress(address string) error {
	timeOut := 5 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeOut)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
