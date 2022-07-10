package config

import "errors"

var (
	ErrNoProvider = errors.New("no provider")

	ErrNoCodec = errors.New("no codec")

	ErrNoName = errors.New("no name")

	ErrConfigNotExist = errors.New("config not exist")
)

type Config interface {
	Get(string, interface{}) error
	IsSet(string) bool
	GetProvider() Provider
	GetInt(string, int) int
	GetInt32(string, int32) int32
	GetInt64(string, int64) int64
	GetUint(string, uint) uint
	GetUint32(string, uint32) uint32
	GetUint64(string, uint64) uint64
	GetFloat32(string, float32) float32
	GetFloat64(string, float64) float64
	GetString(string, string) string
	GetBool(string, bool) bool
}

type ConfigOption func(*ConfigOptions)

type ConfigOptions struct {
	Name     string
	Provider Provider
	Codec    Codec
}

func WithName(name string) ConfigOption {
	return func(opts *ConfigOptions) {
		opts.Name = name
	}
}

func WithProvider(provider Provider) ConfigOption {
	return func(opts *ConfigOptions) {
		opts.Provider = provider
	}
}

func WithCodec(codec Codec) ConfigOption {
	return func(opts *ConfigOptions) {
		opts.Codec = codec
	}
}

type Loader interface {
	Load(opts ...ConfigOption) (Config, error)
	Reload(opts ...ConfigOption) error
}

type ProviderCallback func(string, []byte)

type Provider interface {
	Name() string
	Load() error
	Reload() error
	Read(string) ([]byte, error)
	Watch(ProviderCallback)
}

var providerMap = make(map[string]Provider)

func RegisterProvider(name string, provider Provider) {
	providerMap[name] = provider
}

func GetProvider(name string) Provider {
	return providerMap[name]
}

type Codec interface {
	Name() string
	Unmarshal([]byte, interface{}) error
}

var codecMap = make(map[string]Codec)

func RegisterCodec(name string, codec Codec) {
	codecMap[name] = codec
}

func GetCodec(name string) Codec {
	return codecMap[name]
}

func Load(loader Loader, opts ...ConfigOption) (Config, error) {
	return loader.Load(opts...)
}

func Reload(loader Loader, opts ...ConfigOption) error {
	return loader.Reload(opts...)
}
