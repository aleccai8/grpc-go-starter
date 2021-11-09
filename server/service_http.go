package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewHttpService(opts ...Option) Service {
	gOption := &Options{}
	for _, f := range opts {
		f(gOption)
	}
	return &HttpService{
		opts: gOption,
	}
}

func NewHttpServiceDesc(registrar HttpRegistrar) *ServiceDescHTTP {
	return &ServiceDescHTTP{registrar: registrar}
}

type HttpService struct {
	desc     *ServiceDescHTTP
	server   *http.Server
	serveMux *runtime.ServeMux
	dialConn *grpc.ClientConn
	opts     *Options
}

func (s *HttpService) Register(serviceDesc interface{}, nil interface{}) {
	desc, ok := serviceDesc.(*ServiceDescHTTP)
	if !ok {
		fmt.Println(errors.New("service desc type invalid"))
		return
	}
	s.desc = desc
	filters := assertHttpOptions(s.opts.Filters)
	opts := assertHttpOptions(s.opts.Others)
	opts = append(opts, filters...)
	s.serveMux = runtime.NewServeMux(opts...)
}

func (s *HttpService) Serve() error {
	conn, err := grpc.DialContext(
		context.Background(),
		s.opts.Target,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	s.dialConn = conn
	err = s.desc.registrar(context.Background(), s.serveMux, s.dialConn)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", s.opts.Port),
		Handler: s.serveMux,
	}
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			return
		}
	}()
	err = s.opts.Registry.Register(s.opts.ServiceName)
	return nil
}

func (s HttpService) Close(c chan struct{}) error {
	return s.server.Close()
}

type HttpRegistrar func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type ServiceDescHTTP struct {
	registrar HttpRegistrar
}

func assertHttpOptions(inters []interface{}) []runtime.ServeMuxOption {
	opts := make([]runtime.ServeMuxOption, 0)
	for _, inter := range inters {
		opt := inter.(runtime.ServeMuxOption)
		opts = append(opts, opt)
	}
	return opts
}
