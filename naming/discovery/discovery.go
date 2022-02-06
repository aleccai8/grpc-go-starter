package discovery

import (
	"context"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"sync"
)

var (
	discoveries = make(map[string]Discovery)
	lock        = sync.RWMutex{}
)

var (
	ErrDiscoveryNotFound = fmt.Errorf("discovery not found")
)

const (
	PluginType = "discovery"
)

type Options struct {
	Context   context.Context
	Namespace string
}

type Option func(opt *Options)

type Discovery interface {
	List(name string, opts ...Option) ([]*registry.Node, error)

	Target(target string, opts ...Option) (string, error)
}
