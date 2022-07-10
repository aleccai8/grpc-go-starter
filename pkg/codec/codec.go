package codec

import "sync"

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

var (
	clientCodecs = make(map[string]Codec)
	serverCodecs = make(map[string]Codec)
	lock         sync.RWMutex
)

func Register(name string, serverCodec Codec, clientCodec Codec) {
	lock.Lock()
	serverCodecs[name] = serverCodec
	clientCodecs[name] = clientCodec
	lock.Unlock()
}

func GetServer(name string) Codec {
	lock.RLock()
	c := serverCodecs[name]
	lock.RUnlock()
	return c
}

func GetClient(name string) Codec {
	lock.RLock()
	c := clientCodecs[name]
	lock.RUnlock()
	return c
}
