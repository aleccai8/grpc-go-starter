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

// WithAddress 指定server监听地址 ip:port or :port
func WithAddress(s string) Option {
	return func(opts *Options) {
		opts.Address = s
	}
}
