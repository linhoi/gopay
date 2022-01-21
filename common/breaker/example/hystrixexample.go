package example

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func example() {
	var circuitBreakerName string
	hystrix.ConfigureCommand(circuitBreakerName, hystrix.CommandConfig{
		Timeout:                int(5 * time.Second), // 执行command的超时时间为3s
		MaxConcurrentRequests:  100,                  // command的最大并发量
		RequestVolumeThreshold: 100,                  // 统计窗口10s内的请求数量，达到这个请求数量后才去判断是否要开启熔断
		SleepWindow:            int(5 * time.Second), // 当熔断器被打开后，SleepWindow的时间就是控制过多久后去尝试服务是否可用了
		ErrorPercentThreshold:  20,                   // 错误百分比，请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
	})
	hystrix.GetCircuitSettings()

	hystrix.DoC(
		context.Background(),
		circuitBreakerName,
		func(ctx context.Context) error {
			return nil

		},
		func(ctx context.Context,err error) error {
			return nil
		})

	circuitBreaker, _, err := hystrix.GetCircuit(circuitBreakerName)
	if err != nil {
	}

	if circuitBreaker.IsOpen() {

	}
}

