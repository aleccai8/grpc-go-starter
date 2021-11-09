package server

import (
	"github.com/zhengheng7913/grpc-go-starter/config"
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
)

type ServiceConstructor func(starter *ServiceStarter, opts ...Option) Service

type ServiceType string

type ServiceStarter struct {
	Global         *config.Config
	Current        *ServiceConfig
	CurrentDecoder *config.YamlNodeDecoder
}

type ServiceConfig struct {
	Name     string   `yaml:"name"`
	Protocol string   `yaml:"protocol"`
	Port     uint16   `yaml:"port"`
	Registry string   `yaml:"registry"`
	Filters  []string `yaml:"filters"`
}

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{})

	Serve() error

	Close(chan struct{}) error
}

type Options struct {
	Registry registry.Registry

	Filters []filter.Filter

	Customs []interface{}
}

type Option func(*Options)
