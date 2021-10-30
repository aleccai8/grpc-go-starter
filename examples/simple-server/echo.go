package main

import (
	"context"
	"github.com/zhengheng7913/grpc-config/examples/simple-server/proto"
)

var (
	ExampleService     = "grpc.one.user_server.UserService"
	ExampleServiceHTTP = "grpc.one.user_server.UserServiceHTTP"
)

type EchoServiceImpl struct {
	proto.UnimplementedEchoServiceServer
}

func (e EchoServiceImpl) Echo(ctx context.Context, request *proto.EchoRequest) (*proto.EchoReply, error) {

	return &proto.EchoReply{
		Message: "hello world",
	}, nil
}
