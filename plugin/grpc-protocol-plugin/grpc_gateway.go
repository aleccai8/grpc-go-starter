package grpc_protocol_plugin

import (
	"context"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/registry"
	"github.com/zhengheng7913/grpc-go-starter/pkg/server"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-protocol-plugin/utils"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewGatewayService(opts ...server.Option) server.Service {
	gOption := &server.Options{}
	for _, f := range opts {
		f(gOption)
	}
	return &GatewayService{
		opts: gOption,
	}
}

type GatewayService struct {
	factory  interface{}
	server   *http.Server
	serveMux *runtime.ServeMux
	dialConn *grpc.ClientConn
	opts     *server.Options
}

func (s *GatewayService) Register(factory interface{}, _ interface{}) {
	var opts []runtime.ServeMuxOption
	opts = append(opts, utils.ArrayConvert[runtime.ServeMuxOption](s.opts.Filters)...)
	opts = append(opts, utils.ArrayConvert[runtime.ServeMuxOption](s.opts.Others)...)
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

func (s *GatewayService) Serve() error {
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
		Addr:    fmt.Sprintf("%v:%v", s.opts.Host, s.opts.Port),
		Handler: s.serveMux,
	}
	go func() {
		if s.opts.Registry != nil {
			err = s.opts.Registry.Register(
				s.opts.Name,
				registry.WithNamespace(s.opts.Namespace),
				registry.WithHost(s.opts.Host),
				registry.WithProtocol("http"),
				registry.WithServiceName(s.opts.ServiceName),
				registry.WithPort(s.opts.Port),
			)
			if err != nil {
				log.Println(err)
			}
		}

		err = s.server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (s GatewayService) Close() error {
	if s.opts.Registry != nil {
		s.opts.Registry.Deregister(s.opts.Name)
	}
	s.server.Close()
	return nil
}
