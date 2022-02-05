package config

import "os"

func NewOsEnvProvider() ConfigProvider {
	return &OSEnvProvider{}
}

type OSEnvProvider struct {
}

func (provider *OSEnvProvider) Name() string {
	return "os_env"
}

func (provider *OSEnvProvider) Load() error {
	return nil
}

func (provider *OSEnvProvider) Reload() error {
	return nil
}

func (provider *OSEnvProvider) Read(key string) ([]byte, error) {
	return []byte(os.Getenv(key)), nil
}

// 环境变量无法监听修改
func (provider *OSEnvProvider) Watch(callback ProviderCallback) {

}
