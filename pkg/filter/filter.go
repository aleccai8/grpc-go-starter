package filter

import (
	"sync"
)

var (
	serverFilters = make(map[string]Filter)
	clientFilters = make(map[string]Filter)
	lock          = sync.RWMutex{}
)

type Filter = interface{}

type Chain = []Filter

func Register(name string, server Filter, client Filter) {
	lock.Lock()
	serverFilters[name] = server
	clientFilters[name] = client
	lock.Unlock()
}

// GetServer 获取server拦截器
func GetServer(name string) Filter {
	lock.RLock()
	f := serverFilters[name]
	lock.RUnlock()
	return f
}

// GetClient 获取client拦截器
func GetClient(name string) Filter {
	lock.RLock()
	f := clientFilters[name]
	lock.RUnlock()
	return f
}
