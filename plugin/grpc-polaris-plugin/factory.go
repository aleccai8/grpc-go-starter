package grpc_polaris_plugin

import (
	"fmt"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/discovery"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/registry"
	"github.com/zhengheng7913/grpc-go-starter/pkg/plugin"
	"google.golang.org/grpc/resolver"
	"io/ioutil"
	"strings"
)

const (
	PluginName = "polaris"
)

var (
	ErrPluginDecoderEmpty = fmt.Errorf("plugin decoder is empty")
)

func init() {
	plugin.Register(PluginName, &Factory{})
}

type Factory struct {
}

func (f *Factory) Destroy() error {
	return nil
}

func (f *Factory) Setup(name string, dec plugin.Decoder) error {
	if dec == nil {
		return ErrPluginDecoderEmpty
	}
	conf := &FactoryConfig{}
	if err := dec.Decode(conf); err != nil {
		return err
	}
	return register(conf)
}

func loadConfiguration(conf *FactoryConfig) (config.Configuration, error) {
	buf, err := ioutil.ReadFile("./polaris.yaml")
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadConfiguration(buf)
	if err != nil {
		return nil, err
	}
	if len(conf.AddressList) > 0 {
		addressList := strings.Split(conf.AddressList, ",")
		cfg.GetGlobal().GetServerConnector().SetAddresses(addressList)
	}
	return cfg, nil
}

func register(conf *FactoryConfig) (err error) {
	cfg, err := loadConfiguration(conf)
	if err != nil {
		return err
	}
	var provider api.ProviderAPI
	var consumer api.ConsumerAPI
	if len(conf.Services) > 0 {
		if provider, err = api.NewProviderAPIByConfig(cfg); err != nil {
			return err
		}
	}
	if len(conf.Clients) > 0 {
		if consumer, err = api.NewConsumerAPIByConfig(cfg); err != nil {
			return err
		}
		resolver.Register(NewPolarisResolverBuilder(consumer))
	}

	for _, client := range conf.Clients {
		dc := &DiscoveryConfig{
			Name:     client.Name,
			MetaData: client.MetaData,
		}
		d := newDiscovery(consumer, dc)
		discovery.Register(dc.Name, d)
	}
	for _, service := range conf.Services {
		rc := &RegistryConfig{
			EnableRegister:     service.EnableRegister,
			HeartBeat:          conf.HeartbeatInterval / 1000,
			ServiceToken:       service.Token,
			DisableHealthCheck: conf.DisableHealthCheck,
			Version:            service.Version,
		}
		reg := newRegistry(provider, rc)
		registry.Register(service.Name, reg)
	}
	return nil
}
