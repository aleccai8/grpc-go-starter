package utils

import (
	"context"
	"github.com/zhengheng7913/grpc-go-starter/pkg/server"
	"google.golang.org/grpc"
)

func ArrayConvert[T any](array []interface{}) []T {
	opts := make([]T, len(array))
	for i, inter := range array {
		opt := inter.(T)
		opts[i] = opt
	}
	return opts
}

func GetContextValueInterceptor(opts *server.Options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = context.WithValue(ctx, "service", opts.ServiceName)
		return handler(ctx, req)
	}
}
