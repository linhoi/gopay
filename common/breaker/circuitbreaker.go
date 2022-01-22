package breaker

import (
	"context"
	"errors"
)

type BreakerType string

const (
	Hystrix BreakerType = "hystrix"
)

type runFuncC func(context.Context) error
type fallbackFuncC func(context.Context, error) error

// ICircuitBreaker ...
type ICircuitBreaker interface {
	IsOpen() bool

	DoC(ctx context.Context, run func(context.Context) error, fallback func(context.Context, error) error) error

	Do(run func() error, fallback func(error) error) error

	GetName() string
}

func New(breakerType BreakerType, c Config) (ICircuitBreaker, error) {

	switch breakerType {
	case Hystrix:
		return NewHystrixBreaker(&c)
	default:
		return nil, errors.New("unkonwn breaker type")

	}
}
