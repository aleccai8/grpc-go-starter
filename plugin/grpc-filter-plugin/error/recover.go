package error

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

const (
	GrpcRecoverServerFilter = "grpc-recover"
)

func UnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			grpclog.Errorf("panic recover: %v", e)
			resp = nil
			err = status.Errorf(codes.Internal, "server error: %s", e)
		}
	}()
	return handler(ctx, req)
}
