package client

import (
	"context"
	"github.com/zhengheng7913/grpc-go-starter/pkg/filter"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/discovery"
)

var (
	implementMap = make(map[string]func(opts ...Option) Client[interface{}])
)

func Get(name string) func(opt ...Option) Client[interface{}] {
	return implementMap[name]
}

func Register(name string, cons func(opts ...Option) Client[interface{}]) {
	implementMap[name] = cons
}

func WithTarget(name string) Option {
	return func(opt *Options) {
		opt.Target = name
	}
}

func WithServiceName(name string) Option {
	return func(opt *Options) {
		opt.ServiceName = name
	}
}

func WithNamespace(namespace string) Option {
	return func(opt *Options) {
		opt.Namespace = namespace
	}
}

func WithDiscovery(d discovery.Discovery) Option {
	return func(opt *Options) {
		opt.Discovery = d
	}
}

func WithFilter(filters []filter.Filter) Option {
	return func(opt *Options) {
		opt.Filters = filters
	}
}

func WithProtocol(protocol string) Option {
	return func(opt *Options) {
		opt.Protocol = protocol
	}
}

func WithSrcService(name string) Option {
	return func(opt *Options) {
		opt.SrcServiceName = name
	}
}

type Proxy[T interface{}] struct {
	c Client[interface{}]
}

func (p *Proxy[T]) RealClient(ctx context.Context) T {
	return p.c.RealClient(ctx).(T)
}

func (p *Proxy[T]) Register(realClient any, opts ...Option) {
	p.c.Register(realClient, opts...)
}

func UseClient[T interface{}](c Client[interface{}]) Client[T] {
	return &Proxy[T]{
		c: c,
	}
}

type Options struct {
	Discovery      discovery.Discovery
	Namespace      string
	Target         string
	Protocol       string
	ServiceName    string
	SrcServiceName string
	Metadata       map[string]string
	Filters        []filter.Filter
}

type Option func(opt *Options)

type Client[T interface{}] interface {
	RealClient(ctx context.Context) T

	Register(realClient any, opts ...Option)
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[string]Client[interface{}]),
	}
}

type Clients struct {
	clients map[string]Client[interface{}]
}

func (m *Clients) AddClient(name string, client Client[interface{}]) {
	m.clients[name] = client
}

func (m *Clients) Client(name string) Client[interface{}] {
	return m.clients[name]
}
