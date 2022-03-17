package server

import (
	"context"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
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

type HttpService struct {
	factory  interface{}
	server   *http.Server
	serveMux *runtime.ServeMux
	dialConn *grpc.ClientConn
	opts     *Options
}

func (s *HttpService) Register(factory interface{}, _ interface{}) {
	var opts []runtime.ServeMuxOption
	opts = append(opts, arrayConvert[runtime.ServeMuxOption](s.opts.Filters)...)
	opts = append(opts, arrayConvert[runtime.ServeMuxOption](s.opts.Others)...)
	opts = append(opts,
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}}),
	)
	s.serveMux = runtime.NewServeMux(opts...)
	s.factory = factory
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
	factory := s.factory.(func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error)
	err = factory(context.Background(), s.serveMux, s.dialConn)
	if err != nil {
		return err
	}
	s.server = &http.Server{
		// TODO: 添加ipport支持
		Addr:    fmt.Sprintf("%v:%v", s.opts.Host, s.opts.Port),
		Handler: s.serveMux,
	}
	go func() {
		defer s.opts.Registry.Deregister(s.opts.ServiceName)
		err = s.opts.Registry.Register(
			s.opts.ServiceName,
			registry.WithNamespace(s.opts.Namespace),
			registry.WithHost(s.opts.Host),
			registry.WithProtocol(ProtocolNameHTTP),
			registry.WithServiceName(s.opts.ServiceName),
			registry.WithPort(s.opts.Port),
		)
		if err != nil {
			grpclog.Errorln(err)
		}
		err := s.server.ListenAndServe()
		if err != nil {
			grpclog.Fatalln(err)
		}
	}()

	return nil
}

func (s HttpService) Close(c chan struct{}) error {
	return s.server.Close()
}
