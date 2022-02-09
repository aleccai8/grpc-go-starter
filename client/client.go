package client

import (
	"context"
	"errors"
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/discovery"
)

var (
	implementMap = make(map[string]func(opts ...Option) Client)
)

var (
	ErrClientInvalid = errors.New("err client invalid")
)

func init() {
	implementMap[GrpcProtocol] = NewGrpcClient
	//implementMap[HttpProtocol] = NewHttpClient
}

func Get(name string) func(opt ...Option) Client {
	return implementMap[name]
}

func Invoke[T1 any, T2 any](client Client, context context.Context, method any, req T1, opts ...Option) (reply T2, err error) {
	options := applyOption(opts...)
	if client == nil {
		return reply, ErrClientInvalid
	}
	r, err := client.Invoke(context, method, req, options)
	reply = r.(T2)
	return reply, err
}

func Register[T any](client Client, realClient T, opts ...Option) {
	options := applyOption(opts...)
	if client == nil {
		panic(ErrClientInvalid)
	}
	client.Register(realClient, options)
}

func RealClient[T any](client Client) T {
	if client == nil {
		panic(ErrClientInvalid)
	}
	return client.RealClient().(T)
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

type Options struct {
	Discovery   discovery.Discovery
	ServiceName string
	Namespace   string
	Filters     []filter.Filter
}

type Option func(opt *Options)

type Client interface {
	Invoke(context context.Context, method any, req any, options *Options) (any, error)

	RealClient() any

	Register(realClient any, options *Options)
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[string]Client),
	}
}

type Clients struct {
	clients map[string]Client
}

func (m *Clients) AddClient(name string, client Client) {
	m.clients[name] = client
}

func (m *Clients) Client(name string) Client {
	return m.clients[name]
}
