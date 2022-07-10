package server

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{})

	Serve() error

	Close() error
}
