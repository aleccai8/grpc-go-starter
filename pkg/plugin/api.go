package plugin

import "time"

type Factory interface {
	Setup(name string, dec Decoder) error

	Destroy() error
}

type Decoder interface {
	Decode(cfg interface{}) error
}

func Register(name string, f Factory) {
	factories[name] = f
}

func Get(name string) Factory {
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
