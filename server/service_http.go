package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type HttpServiceConfig struct {
	ServiceConfig
	Target string `yaml:"target"`
}

type HttpService struct {
	desc     *ServiceDescHTTP
	cfg      *HttpServiceConfig
	server   *http.Server
	serveMux *runtime.ServeMux
	dialConn *grpc.ClientConn
	opt      *Options
}

func (s *HttpService) Register(serviceDesc interface{}, nil interface{}) {
	desc, ok := serviceDesc.(*ServiceDescHTTP)
	if !ok {
		fmt.Println(errors.New("service desc type invalid"))
		return
	}
	s.desc = desc
	filters, _ := assertHttpOptions(s.opt.Filters)
	opts, _ := assertHttpOptions(s.opt.Customs)
	opts = append(opts, filters...)
	s.serveMux = runtime.NewServeMux(opts...)
}

func (s *HttpService) Serve() error {
	conn, err := grpc.DialContext(
		context.Background(),
		s.cfg.Target,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	s.dialConn = conn
	err = s.desc.registrar(context.Background(), s.serveMux, s.dialConn)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", s.cfg.Port),
		Handler: s.serveMux,
	}
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			return
		}
	}()
	err = s.opt.Registry.Register(s.cfg.Name)
	return nil
}

func (s HttpService) Close(c chan struct{}) error {
	return s.server.Close()
}

type HttpRegistrar func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type ServiceDescHTTP struct {
	registrar HttpRegistrar
}

func assertHttpOptions(inters ...interface{}) ([]runtime.ServeMuxOption, bool) {
	opts := make([]runtime.ServeMuxOption, 0)
	for _, inter := range inters {
		opt, ok := inter.(runtime.ServeMuxOption)
		if !ok {
			return nil, false
		}
		opts = append(opts, opt)
	}
	return opts, true
}

func dessertHttpOptions(opts ...runtime.ServeMuxOption) []interface{} {
	inters := make([]interface{}, len(opts))
	for _, opt := range opts {
		inters = append(inters, opt)
	}
	return inters
}
