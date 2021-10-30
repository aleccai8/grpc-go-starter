package server

import (
	"github.com/zhengheng7913/grpc-config/config"
)

type ServiceConstructor func(cfg *config.ServiceConfig, opts ...Option) Service

type ServiceType string

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{})

	Serve() error

	Close(chan struct{}) error
}

type ServiceOptions interface {
	ProtocolName() string

	Apply(inters ...interface{})
}
