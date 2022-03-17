package server

import (
	"context"
	"google.golang.org/grpc"
)

const (
	KeyService = "service"
)

func arrayConvert[T any](array []interface{}) []T {
	opts := make([]T, len(array))
	for i, inter := range array {
		opt := inter.(T)
		opts[i] = opt
	}
	return opts
}

func GetContextValueInterceptor(opts *Options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = context.WithValue(ctx, KeyService, opts.ServiceName)
		return handler(ctx, req)
	}
}
