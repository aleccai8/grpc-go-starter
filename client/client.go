package client

import "context"

var (
	constructors = make(map[string]func() Client)
)

func init() {
	constructors[GrpcProtocol] = NewGrpcClient
	constructors[HttpProtocol] = NewHttpClient
}

func WithDiscovery(name string) Option {
	return func(opt *Options) {
		opt.Discovery = name
	}
}

func WithName(name string) Option {
	return func(opt *Options) {
		opt.Name = name
	}
}

func WithNamespace(namespace string) Option {
	return func(opt *Options) {
		opt.Namespace = namespace
	}
}

type Options struct {
	Discovery string
	Name      string
	Namespace string
}

type Option func(opt *Options)

type Client interface {
	Invoke(context context.Context, method string, req interface{}, opts ...Option) (interface{}, error)

	Register(realClient interface{}, opts ...Option)
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
