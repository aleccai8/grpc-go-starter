package discovery

import "context"

func Register(name string, discovery Discovery) {
	lock.Lock()
	defer lock.Unlock()
	discoveries[name] = discovery
}

func Get(name string) Discovery {
	lock.RLock()
	defer lock.RUnlock()
	return discoveries[name]
}

func WithNamespace(namespace string) Option {
	return func(opt *Options) {
		opt.Namespace = namespace
	}
}

func WithContext(context context.Context) Option {
	return func(opt *Options) {
		opt.Context = context
	}
}

func WithSrcService(name string) Option {
	return func(opt *Options) {
		opt.SrcService = name
	}
}

func WithProtocol(name string) Option {
	return func(opt *Options) {
		opt.Protocol = name
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(opt *Options) {
		opt.Metadata = metadata
	}
}
