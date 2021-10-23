package grpc

import (
	"github.com/zhengheng7913/grpc-config/config"
	"github.com/zhengheng7913/grpc-config/server"
	"google.golang.org/grpc"
)

type Service struct {
	server *grpc.Server
	cfg    *config.ServiceConfig
	opt    *Option
}

func (g *Service) Register(serviceDesc *server.ServiceDesc, serviceImpl interface{}) error {
	panic("implement me")
}

func (g *Service) Serve() error {
	return nil
}

func (g *Service) Close(c chan struct{}) error {
	panic("implement me")
}
