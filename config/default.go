package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidUnmarshalType = errors.New("invalid unmarshal type")

	ErrInvalid
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
	*(output.(*string)) = string(input)
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

type YAMLCodec struct {
}

func (y *YAMLCodec) Name() string {
	return "yaml"
}

func (y *YAMLCodec) Unmarshal(input []byte, output interface{}) error {
	return yaml.Unmarshal(input, output)
}

func NewDefaultLoader() ConfigLoader {
	return &DefaultLoader{
		configMap: make(map[string]Config),
	}
}

type DefaultLoader struct {
	configMap map[string]Config
	rwl       sync.RWMutex
}

func (loader *DefaultLoader) Load(opts ...ConfigOption) (Config, error) {
	options := loader.applyOptions(opts...)
	if err := loader.checkOptions(options); err != nil {
		return nil, err
	}

	key := loader.getKey(options)
	loader.rwl.RLock()
	if c, ok := loader.configMap[key]; ok {
		loader.rwl.RUnlock()
		return c, nil
	}
	loader.rwl.RUnlock()
	config := NewDefaultConfig(options)
	if err := config.GetProvider().Load(); err != nil {
		return nil, err
	}
	loader.rwl.Lock()
	loader.configMap[key] = config
	loader.rwl.Unlock()
	return config, nil
}

func (loader *DefaultLoader) Reload(opts ...ConfigOption) error {
	options := loader.applyOptions(opts...)
	key := loader.getKey(options)
	loader.rwl.RLock()
	defer loader.rwl.RUnlock()
	if config, ok := loader.configMap[key]; ok {
		return config.GetProvider().Reload()
	}
	return ErrConfigNotExist
}

func (loader *DefaultLoader) checkOptions(options *ConfigOptions) error {
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

func (loader *DefaultLoader) getKey(options *ConfigOptions) string {
	return fmt.Sprintf("%s.%s.%s", options.Codec.Name(), options.Provider.Name(), options.Name)
}

func (loader *DefaultLoader) applyOptions(opts ...ConfigOption) *ConfigOptions {
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
	provider ConfigProvider
	codec    Codec
}

func (config *DefaultConfig) GetProvider() ConfigProvider {
	return config.provider
}

func (config *DefaultConfig) GetCodec() Codec {
	return config.codec
}

func (config *DefaultConfig) findWithDefaultValue(key string, defaultValue interface{}) interface{} {
	b, err := config.GetProvider().Read(key)
	if err != nil {
		return defaultValue
	}
	var v interface{} = nil
	err = config.GetCodec().Unmarshal(b, &v)
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

func (config *DefaultConfig) Get(key string, defaultValue interface{}) interface{} {
	return config.findWithDefaultValue(key, defaultValue)
}

func (config *DefaultConfig) IsSet(_ string) bool {
	panic("not implemented") // TODO: Implement
}

func (config *DefaultConfig) GetInt(key string, defaultValue int) int {
	return cast.ToInt(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetInt32(key string, defaultValue int32) int32 {
	return cast.ToInt32(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetInt64(key string, defaultValue int64) int64 {
	return cast.ToInt64(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetUint(key string, defaultValue uint) uint {
	return cast.ToUint(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetUint32(key string, defaultValue uint32) uint32 {
	return cast.ToUint32(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetUint64(key string, defaultValue uint64) uint64 {
	return cast.ToUint64(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetFloat32(key string, defaultValue float32) float32 {
	return cast.ToFloat32(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetFloat64(key string, defaultValue float64) float64 {
	return cast.ToFloat64(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetString(key string, defaultValue string) string {
	return cast.ToString(config.findWithDefaultValue(key, defaultValue))
}

func (config *DefaultConfig) GetBool(key string, defaultValue bool) bool {
	return cast.ToBool(config.findWithDefaultValue(key, defaultValue))
}
