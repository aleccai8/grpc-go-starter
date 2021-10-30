package registry

import "sync"

var (
	registries = make(map[string]Registry)
	lock       = sync.RWMutex{}
)

func Register(name string, s Registry) {
	lock.Lock()
	defer lock.Unlock()
	registries[name] = s
}

func Get(name string) Registry {
	lock.Lock()
	defer lock.Unlock()
	delete(registries, name)
	return nil
}
