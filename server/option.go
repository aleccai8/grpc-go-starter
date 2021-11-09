package server

import (
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
)

type Options struct {
	Registry registry.Registry

	ServiceOptions ServiceOptions
}

type Option func(*Options)
