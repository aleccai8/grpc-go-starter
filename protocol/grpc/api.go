package grpc

import (
	"github.com/zhengheng7913/grpc-config/server"
)

const (
	ProtocolName = "grpc"
)

func NewGrpcService(opts ...server.Options) server.Service {
	// 初始化GrpcServer
	gOption := &Option{}
	for _, f := range opts {
		f(gOption)
	}
	return &Service{
		opt: gOption,
	}
}
