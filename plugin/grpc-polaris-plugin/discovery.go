package grpc_polaris_plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/polarismesh/polaris-go/api"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/discovery"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/registry"
)

func NewDiscovery(consumer api.ConsumerAPI, cfg *DiscoveryConfig) discovery.Discovery {
	return newDiscovery(consumer, cfg)
}

func newDiscovery(consumer api.ConsumerAPI, cfg *DiscoveryConfig) *Discovery {
	return &Discovery{
		consumer: consumer,
		cfg:      cfg,
	}
}

type Discovery struct {
	consumer api.ConsumerAPI
	cfg      *DiscoveryConfig
}

func (d *Discovery) List(name string, opts ...discovery.Option) ([]*registry.Node, error) {
	return nil, nil
}

func (d *Discovery) Target(target string, opts ...discovery.Option) (string, error) {
	options := &discovery.Options{}
	for _, o := range opts {
		o(options)
	}
	if options.Metadata == nil {
		options.Metadata = make(map[string]string)
	}
	if options.Protocol != "" {
		options.Metadata["protocol"] = options.Protocol
	}
	dialOptions := &DialOptions{
		Namespace:   options.Namespace,
		SrcService:  options.SrcService,
		SrcMetadata: d.cfg.MetaData,
		DstMetadata: options.Metadata,
	}
	str, err := json.Marshal(dialOptions)
	if err != nil {
		return "", fmt.Errorf("marshal dialOptions error: %s", err)
	}
	endpoint := base64.URLEncoding.EncodeToString(str)
	return scheme + "://" + target + "/" + endpoint, nil
}
