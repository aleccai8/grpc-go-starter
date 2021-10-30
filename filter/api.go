package filter

import (
	"github.com/zhengheng7913/grpc-config/server"
	"sync"
)

var (
	chains = make(map[string][]server.Options)
	lock   = sync.RWMutex{}
)

func Register(serviceName string, chain server.Options) {
	lock.Lock()
	defer lock.Unlock()
	if chains[serviceName] == nil {
		chains[serviceName] = make([]server.Options, 0)
	}
	chains[serviceName] = append(chains[serviceName], chain)
}

func Get(name string) []server.Options {
	lock.RLock()
	defer lock.RUnlock()
	return chains[name]
}
