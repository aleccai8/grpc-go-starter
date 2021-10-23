package server

type ServiceConstructor func(opts ...Options) Service

type ServiceType string

type Service interface {
	Register(serviceDesc *ServiceDesc, serviceImpl interface{}) error

	Serve() error

	Close(chan struct{}) error
}

const (
	ServiceMethodGrpc = "grpc"
	ServiceMethodHttp = "http"
)

type ServiceDesc struct {
	ServiceName string
	Method      ServiceType
}
