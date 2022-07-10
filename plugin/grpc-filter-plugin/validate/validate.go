package validate

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GrpcUnaryServerFilter = "grpc-unary-validate"
)

type validator interface {
	Validate() error

	ValidateAll() error
}

func validate(entity interface{}) error {
	v, ok := entity.(validator)
	if !ok {
		return nil
	}
	return v.ValidateAll()
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if err := validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "client error: %v", err.Error())
	}
	reply, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	if err := validate(reply); err != nil {
		return nil, status.Errorf(codes.Internal, "server error: %v", err.Error())
	}
	return reply, nil
}
