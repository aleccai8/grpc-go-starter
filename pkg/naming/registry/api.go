package registry

func Register(name string, s Registry) {
	lock.Lock()
	defer lock.Unlock()
	registries[name] = s
}

func Get(name string) Registry {
	lock.RLock()
	defer lock.RUnlock()
	return registries[name]
}

func WithNamespace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
	}
}

func WithServiceName(name string) Option {
	return func(options *Options) {
		options.ServiceName = name
	}
}

func WithPort(port uint16) Option {
	return func(options *Options) {
		options.Port = port
	}
}

func WithProtocol(protocol string) Option {
	return func(options *Options) {
		options.Protocol = protocol
	}
}

func WithHost(host string) Option {
	return func(options *Options) {
		options.Host = host
	}
}
