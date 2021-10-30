package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhengheng7913/grpc-config/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

type ServiceHTTP struct {
	desc     *ServiceDescHTTP
	cfg      *config.ServiceConfig
	server   *http.Server
	serveMux *runtime.ServeMux
	dialConn *grpc.ClientConn
}

func (s *ServiceHTTP) Register(serviceDesc interface{}, nil interface{}) {
	desc, ok := serviceDesc.(*ServiceDescHTTP)
	if !ok {
		fmt.Println(errors.New("service desc type invalid"))
		return
	}
	s.desc = desc
}

func (s *ServiceHTTP) Serve() error {
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8000",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	s.dialConn = conn
	s.serveMux = runtime.NewServeMux()
	err = s.desc.registrar(context.Background(), s.serveMux, s.dialConn)
	s.server = &http.Server{
		Addr:    ":8090",
		Handler: s.serveMux,
	}
	return s.server.ListenAndServe()
}

func (s ServiceHTTP) Close(c chan struct{}) error {
	panic("implement me")
}

type RegistrarHTTP func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type ServiceDescHTTP struct {
	registrar RegistrarHTTP
}

type OptionHTTP struct {
	serviceName    string
	serveMuxOption runtime.ServeMuxOption
}

func (o OptionHTTP) Apply(arg ...interface{}) {
	panic("implement me")
}

func (o OptionHTTP) ServiceName() string {
	return o.serviceName
}

func (o OptionHTTP) ProtocolName() string {
	return ProtocolNameHTTP
}

func assertOptions(inters ...interface{}) ([]Option, bool) {
	opts := make([]Option, len(inters))
	for _, inter := range inters {
		opt, ok := inter.(Option)
		if !ok {
			return nil, false
		}
		opts = append(opts, opt)
	}
	return opts, true
}

func dessertOptions(opts ...runtime.ServeMuxOption) []interface{} {
	inters := make([]interface{}, len(opts))
	for _, opt := range opts {
		inters = append(inters, opt)
	}
	return inters
}

const (
	HeaderErrorMsg  = "grpc-error-msg"
	HeaderErrorCode = "grpc-error-code"
)

type ErrorHandlerHTTP struct {
}

func (h ErrorHandlerHTTP) Apply(arg ...interface{}) {
	panic("implement me")
}

func (h ErrorHandlerHTTP) ProtocolName() string {
	return "http"
}

func httpErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	w.Header().Set(HeaderErrorMsg, s.Message())
	w.Header().Set(HeaderErrorCode, strconv.Itoa(int(s.Code())))

	st := runtime.HTTPStatusFromCode(s.Code())
	if customStatus != nil {
		st = customStatus.HTTPStatus
	}

	w.WriteHeader(st)
}
