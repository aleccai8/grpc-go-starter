package plugin

import "time"

type Factory interface {
	Type() string

	Setup(name string, dec Decoder) error

	Destroy() error
}

type Decoder interface {
	Decode(cfg interface{}) error
}

func Register(name string, f Factory) {
	factories, ok := plugins[f.Type()]
	if !ok {
		plugins[f.Type()] = map[string]Factory{
			name: f,
		}
		return
	}
	factories[name] = f
}

func Get(typ string, name string) Factory {
	factories, ok := plugins[typ]
	if !ok {
		return nil
	}
	return factories[name]
}

// WaitForDone 挂住等待所有插件初始化完成，可自己设置超时时间。
func WaitForDone(timeout time.Duration) bool {
	select {
	case <-done:
		return true
	case <-time.After(timeout):
	}
	return false
}
