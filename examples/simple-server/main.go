package main

import (
	"github.com/zhengheng7913/grpc-config"
	"github.com/zhengheng7913/grpc-config/examples/simple-server/proto"
	"github.com/zhengheng7913/grpc-config/server"
)

func main() {
	s := grpc_config.NewServer()

	proto.RegisterEchoServiceServer(
		server.WithServiceRegisterAdapter(s.Service(ExampleService)),
		&EchoServiceImpl{},
	)

	s.Service(ExampleServiceHTTP).Register(
		server.NewHttpServiceDesc(proto.RegisterEchoServiceHandler),
		nil,
	)

	s.Serve()
}
