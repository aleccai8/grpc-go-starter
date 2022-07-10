package utils

import "github.com/zhengheng7913/grpc-go-starter/pkg/client"

func ApplyOption(opts ...client.Option) *client.Options {
	options := &client.Options{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
