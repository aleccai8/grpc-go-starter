package server

type ServiceConstructor func(opts ...Option) Service

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{})

	Serve() error

	Close(chan struct{}) error
}
