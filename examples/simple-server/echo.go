package main

import (
	"context"
	"github.com/zhengheng7913/grpc-config/examples/simple-server/proto"
)

var (
	ExampleService     = "grpc.examples.simple-server.ExampleService"
	ExampleServiceHTTP = "grpc.examples.simple-server.ExampleServiceHTTP"
)

type EchoServiceImpl struct {
	proto.UnimplementedEchoServiceServer
}

func (e EchoServiceImpl) Echo(ctx context.Context, request *proto.EchoRequest) (*proto.EchoReply, error) {

	return &proto.EchoReply{
		Message: "hello world",
	}, nil
}
