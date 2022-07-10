package config

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
	"reflect"
	"sync"
)

var (
	ErrInvalidUnmarshalType = errors.New("invalid unmarshal type")
)

func NewKVCodec() Codec {
	return &KVCodec{}
}

type KVCodec struct {
}

func (k *KVCodec) Name() string {
	return "kv"
}

func (k *KVCodec) Unmarshal(input []byte, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidUnmarshalType
	}
	rv.Elem().Set(reflect.ValueOf(string(input)))
	return nil
}

type JSONCodec struct {
}

func (j *JSONCodec) Name() string {
	return "json"
}

func (j *JSONCodec) Unmarshal(input []byte, output interface{}) error {
	return json.Unmarshal(input, output)
}

func NewYAMLCodec() Codec {
	return &YAMLCodec{}
}

type YAMLCodec struct {
}

func (y *YAMLCodec) Name() string {
	return "yaml"
}

func (y *YAMLCodec) Unmarshal(input []byte, output interface{}) error {
	return yaml.Unmarshal(input, output)
}

func NewDefaultLoader() Loader {
	return &DefaultLoader{
		configMap: make(map[string]Config),
	}
}

type DefaultLoader struct {
	configMap map[string]Config
	rwl       sync.RWMutex
}

func (l *DefaultLoader) Load(opts ...ConfigOption) (Config, error) {
	options := l.applyOptions(opts...)
	if err := l.checkOptions(options); err != nil {
		return nil, err
	}

	key := options.Name
	l.rwl.RLock()
	if c, ok := l.configMap[key]; ok {
		l.rwl.RUnlock()
		return c, nil
	}
	l.rwl.RUnlock()
	config := NewDefaultConfig(options)
	if err := config.GetProvider().Load(); err != nil {
		return nil, err
	}
	l.rwl.Lock()
	l.configMap[key] = config
	l.rwl.Unlock()
	return config, nil
}

func (l *DefaultLoader) Reload(opts ...ConfigOption) error {
	options := l.applyOptions(opts...)
	key := options.Name
	l.rwl.RLock()
	defer l.rwl.RUnlock()
	if config, ok := l.configMap[key]; ok {
		return config.GetProvider().Reload()
	}
	return ErrConfigNotExist
}

func (l *DefaultLoader) checkOptions(options *ConfigOptions) error {
	if options.Name == "" {
		return ErrNoName
	}
	if options.Provider == nil {
		return ErrNoProvider
	}
	if options.Codec == nil {
		return ErrNoCodec
	}
	return nil
}

func (l *DefaultLoader) applyOptions(opts ...ConfigOption) *ConfigOptions {
	options := &ConfigOptions{}
	for _, o := range opts {
		o(options)
	}
	return options
}

func NewDefaultConfig(opts *ConfigOptions) Config {
	return &DefaultConfig{
		provider: opts.Provider,
		codec:    opts.Codec,
	}
}

type DefaultConfig struct {
	provider Provider
	codec    Codec
}

func (c *DefaultConfig) GetProvider() Provider {
	return c.provider
}

func (c *DefaultConfig) GetCodec() Codec {
	return c.codec
}

func (c *DefaultConfig) findWithDefaultValue(key string, defaultValue interface{}) interface{} {
	b, err := c.GetProvider().Read(key)
	if err != nil {
		return defaultValue
	}
	var v interface{} = nil
	err = c.GetCodec().Unmarshal(b, &v)
	if err != nil || v == nil {
		return defaultValue
	}
	switch defaultValue.(type) {
	case bool:
		v, err = cast.ToBoolE(v)
	case string:
		v, err = cast.ToStringE(v)
	case int:
		v, err = cast.ToIntE(v)
	case int32:
		v, err = cast.ToInt32E(v)
	case int64:
		v, err = cast.ToInt64E(v)
	case uint:
		v, err = cast.ToUintE(v)
	case uint32:
		v, err = cast.ToUint32E(v)
	case uint64:
		v, err = cast.ToUint64E(v)
	case float64:
		v, err = cast.ToFloat64E(v)
	case float32:
		v, err = cast.ToFloat32E(v)
	default:
	}

	if err != nil {
		return defaultValue
	}
	return v
}

func (c *DefaultConfig) Get(key string, out interface{}) error {
	b, err := c.GetProvider().Read(key)
	if err != nil {
		return err
	}
	return c.GetCodec().Unmarshal(b, out)
}

func (c *DefaultConfig) IsSet(_ string) bool {
	panic("not implemented") // TODO: Implement
}

func (c *DefaultConfig) GetInt(key string, defaultValue int) int {
	return cast.ToInt(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetInt32(key string, defaultValue int32) int32 {
	return cast.ToInt32(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetInt64(key string, defaultValue int64) int64 {
	return cast.ToInt64(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetUint(key string, defaultValue uint) uint {
	return cast.ToUint(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetUint32(key string, defaultValue uint32) uint32 {
	return cast.ToUint32(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetUint64(key string, defaultValue uint64) uint64 {
	return cast.ToUint64(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetFloat32(key string, defaultValue float32) float32 {
	return cast.ToFloat32(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetFloat64(key string, defaultValue float64) float64 {
	return cast.ToFloat64(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetString(key string, defaultValue string) string {
	return cast.ToString(c.findWithDefaultValue(key, defaultValue))
}

func (c *DefaultConfig) GetBool(key string, defaultValue bool) bool {
	return cast.ToBool(c.findWithDefaultValue(key, defaultValue))
}
