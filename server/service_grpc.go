package server

import (
	"errors"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"google.golang.org/grpc/grpclog"
	"net"

	"google.golang.org/grpc"
)

func NewGrpcService(opts ...Option) Service {
	// 初始化GrpcServer
	gOption := &Options{}
	for _, f := range opts {
		f(gOption)
	}
	return &GrpcService{
		opts: gOption,
	}
}

func newServiceRegisterAdapter(srv Service) grpc.ServiceRegistrar {
	return &ServiceRegisterAdapter{
		service: srv,
	}
}

type GrpcService struct {
	server *grpc.Server
	opts   *Options
}

func (g *GrpcService) Register(serviceDesc interface{}, serviceImpl interface{}) {
	desc, ok := serviceDesc.(*grpc.ServiceDesc)
	if !ok {
		fmt.Println(errors.New("service desc type invalid"))
		return
	}
	filters, _ := assertGrpcOptions(g.opts.Filters)
	opts, _ := assertGrpcOptions(g.opts.Others)
	opts = append(opts, filters...)
	g.server = grpc.NewServer(opts...)
	g.server.RegisterService(desc, serviceImpl)
}

func (g *GrpcService) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v", g.opts.Address))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v ", err)
	}
	go func() {
		defer g.opts.Registry.Deregister(g.opts.ServiceName)
		err = g.opts.Registry.Register(g.opts.ServiceName, registry.WithAddress(g.opts.Address))
		if err != nil {
			grpclog.Errorln(err)
		}
		err := g.server.Serve(lis)
		if err != nil {
			grpclog.Fatalln(err)
		}
	}()

	return nil
}

func (g *GrpcService) Close(c chan struct{}) error {
	g.opts.Registry.Deregister(g.opts.ServiceName)
	g.server.Stop()
	return nil
}

type ServiceRegisterAdapter struct {
	service Service
}

func (s *ServiceRegisterAdapter) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.service.Register(desc, impl)
}

func assertGrpcOptions(inters ...interface{}) ([]grpc.ServerOption, bool) {
	opts := make([]grpc.ServerOption, len(inters))
	for _, inter := range inters {
		opt, ok := inter.(grpc.ServerOption)
		if !ok {
			return nil, false
		}
		opts = append(opts, opt)
	}
	return opts, true
}

func dessertGrpcOptions(opts ...grpc.ServerOption) []interface{} {
	inters := make([]interface{}, len(opts))
	for _, opt := range opts {
		inter := opt.(interface{})
		inters = append(inters, inter)
	}
	return inters
}
