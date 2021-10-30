package server

import (
	"github.com/zhengheng7913/grpc-config/config"
	"github.com/zhengheng7913/grpc-config/naming/registry"
)

type Options struct {
	ServiceConfig *config.ServiceConfig

	Registry registry.Registry

	ServiceOptions ServiceOptions
}

type Option func(*Options)
