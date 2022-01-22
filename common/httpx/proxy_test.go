package httpx

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/linhoi/gopay/common/balance"
	"github.com/linhoi/gopay/common/breaker"
	"github.com/linhoi/gopay/common/netx"
	"github.com/pkg/errors"
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

	err = netx.ValidateUrl("http://127.0.0.1:1087")
	if err != nil {
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
			if res != nil && res.StatusCode != http.StatusOK {
				t.Error("res.StatusCode != http.StatusOK")
			}
			t.Log(res)
		})
	}
}

var aGoogleUrl = `https://www.google.com/robots.txt`

func TestProxyOnProxyChange(t *testing.T) {
	type args struct {
		proxy string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "a in valid proxy",
			args:    args{proxy: "http://localhost:8080"},
			wantErr: false,
		},
		{
			name:    "a valid proxy",
			args:    args{proxy: "http://localhost:108"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := genTestProxyClient()
			if err != nil {
				return
			}
			err = netx.ValidateUrl("http://127.0.0.1:1087")
			if err != nil {
				return
			}

			res, err := p.Get(aGoogleUrl)
			if err == io.EOF {
				return
			}
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

func genTestProxyClient() (*Proxy, error) {
	return NewProxyClient(ProxyConfig{
		ProxyUrl: "http://127.0.0.1:1087",
		BreakerConfig: &breaker.Config{
			Name:                   "proxy breaker",
			Timeout:                int(5 * time.Second), // 执行command的超时时间为3s
			MaxConcurrentRequests:  100,                  // command的最大并发量
			RequestVolumeThreshold: 100,                  // 统计窗口10s内的请求数量，达到这个请求数量后才去判断是否要开启熔断
			SleepWindow:            int(5 * time.Second), // 当熔断器被打开后，SleepWindow的时间就是控制过多久后去尝试服务是否可用了
			ErrorPercentThreshold:  20,                   // 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
		},
	})
}

func genTestProxyClientWithLoadBalance() (*Proxy, error) {
	return NewProxyClient(ProxyConfig{
		ProxyUrl: "http://127.0.0.1:1087",
		BreakerConfig: &breaker.Config{
			Name:                   "proxy breaker",
			Timeout:                int(5 * time.Second), // 执行command的超时时间为3s
			MaxConcurrentRequests:  100,                  // command的最大并发量
			RequestVolumeThreshold: 100,                  // 统计窗口10s内的请求数量，达到这个请求数量后才去判断是否要开启熔断
			SleepWindow:            int(5 * time.Second), // 当熔断器被打开后，SleepWindow的时间就是控制过多久后去尝试服务是否可用了
			ErrorPercentThreshold:  20,                   // 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
		},
		BalanceConfig: &balance.Config{Items: []balance.Item{
			{
				Name:   "row",
				Weight: 1,
			},
			{
				Name:   "proxy",
				Weight: 1,
			},
		}},
	})

}

func BenchmarkProxyWithBalance(b *testing.B) {
	type args struct {
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "a valid proxy",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			p, err := genTestProxyClientWithLoadBalance()
			if err != nil {
				return
			}
			err = netx.ValidateUrl("http://127.0.0.1:1087")
			if err != nil {
				return
			}

			for i := 0; i < b.N; i++ {
				res, err := p.Get(aGoogleUrl)
				if err == io.EOF {
					continue
				}

				if errors.Is(err,context.DeadlineExceeded) {
					b.Log("context dead , pass")
					continue
				}

				if strings.Contains(err.Error(), context.Canceled.Error()) || strings.Contains(err.Error(), context.DeadlineExceeded.Error()){
					b.Log("context error , pass")
					continue
				}
				if err != nil {
					b.Error(err)
				}
				if res == nil {
					continue
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					continue
				}


				if err != nil {
					b.Error(err)
				}
				b.Log(string(body[:10]))
				res.Body.Close()
			}

		})
	}
}
