package breaker

import (
	"context"
	"errors"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
)

type HystrixBreaker struct {
	name    string
	breaker *hystrix.CircuitBreaker
}

func (h HystrixBreaker) Do(run func() error, fallback func(error) error) error {
	return hystrix.Do(
		h.name,
		run,
		fallback)
}

func (h HystrixBreaker) IsOpen() bool {
	return !h.breaker.AllowRequest()
}

func (h HystrixBreaker) DoC(ctx context.Context, run func(context.Context) error, fallback func(context.Context, error) error) error {
	return hystrix.DoC(
		ctx,
		h.name,
		run,
		fallback)
}

func (h HystrixBreaker) GetName() string {
	return h.name
}

// NewHystrixBreaker ...
func NewHystrixBreaker(c *Config) (*HystrixBreaker, error) {
	if c == nil || c.Name == "" {
		return nil, errors.New("breaker name not defined")
	}

	hystrix.ConfigureCommand(c.Name, hystrix.CommandConfig{
		Timeout:                c.Timeout,
		MaxConcurrentRequests:  c.MaxConcurrentRequests,
		RequestVolumeThreshold: c.RequestVolumeThreshold,
		SleepWindow:            c.SleepWindow,
		ErrorPercentThreshold:  c.ErrorPercentThreshold,
	})

	breaker, _, err := hystrix.GetCircuit(c.Name)
	if err != nil {
		fmt.Println("get hystric circuit failed:", err)
		return nil, err
	}

	return &HystrixBreaker{name: c.Name, breaker: breaker}, nil
}

var _ ICircuitBreaker = (*HystrixBreaker)(nil)
