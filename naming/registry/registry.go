package registry

type Registry interface {
	Register(service string, opts ...Option) error
	Deregister(service string) error
}

type Options struct {
	Address string
}

type Option func(*Options)
