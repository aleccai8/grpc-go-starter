package server

import "github.com/zhengheng7913/grpc-config/config"

type ServiceConstructor func(cfg *config.ServiceConfig, opts ...Options) Service

type ServiceType string

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{})

	Serve() error

	Close(chan struct{}) error
}

type Option interface {
	Apply(arg ...interface{})

	ProtocolName() string
}

type Options func(Option)
