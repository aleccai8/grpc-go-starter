package client

import "context"

func applyOption(opts ...Option) *Options {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func RealClient[T any](ctx context.Context, client Client) T {
	return client.RealClient(ctx).(T)
}
