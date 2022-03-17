package registry

import "sync"

var (
	registries = make(map[string]Registry)
	lock       = sync.RWMutex{}
)

const (
	PluginType = "registry"
)

type Registry interface {
	Register(service string, opts ...Option) error
	Deregister(service string) error
}

type Options struct {
	Namespace   string
	ServiceName string
	Host        string
	Port        uint16
	Protocol    string
}

type Option func(*Options)
