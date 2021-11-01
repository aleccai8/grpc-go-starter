package server

import (
	"github.com/zhengheng7913/grpc-config/naming/registry"
)

type Options struct {
	Registry registry.Registry

	ServiceOptions ServiceOptions
}

type Option func(*Options)
