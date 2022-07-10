package grpc_polaris_plugin

import (
	"fmt"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/registry"
	"log"
	"time"
)

const (
	DefaultHeartBeat = 5
	DefaultWeight    = 100
	DefaultTTL       = 5
)

type Registry struct {
	provider api.ProviderAPI
	cfg      *RegistryConfig
}

func NewRegistry(provider api.ProviderAPI, cfg *RegistryConfig) registry.Registry {
	return newRegistry(provider, cfg)
}

func newRegistry(provider api.ProviderAPI, cfg *RegistryConfig) *Registry {
	if cfg.HeartBeat == 0 {
		cfg.HeartBeat = DefaultHeartBeat
	}
	if cfg.Weight == 0 {
		cfg.Weight = DefaultWeight
	}
	if cfg.TTL == 0 {
		cfg.TTL = DefaultTTL
	}
	if cfg.Services == nil {
		cfg.Services = make(map[string]*RegistryService)
	}
	return &Registry{
		provider: provider,
		cfg:      cfg,
	}
}

func (r *Registry) Register(sn string, opt ...registry.Option) error {
	opts := &registry.Options{}
	for _, fn := range opt {
		fn(opts)
	}
	rs := &RegistryService{
		Namespace:   opts.Namespace,
		Port:        int(opts.Port),
		Host:        opts.Host,
		Protocol:    opts.Protocol,
		ServiceName: opts.ServiceName,
	}
	r.cfg.Services[sn] = rs
	if r.cfg.EnableRegister {
		if err := r.register(rs); err != nil {
			return err
		}
	}
	go r.heartBeat(rs)
	return nil
}

func (r *Registry) register(rs *RegistryService) error {
	req := &api.InstanceRegisterRequest{
		InstanceRegisterRequest: model.InstanceRegisterRequest{
			Namespace:    rs.Namespace,
			Service:      rs.ServiceName,
			Host:         rs.Host,
			Port:         rs.Port,
			ServiceToken: r.cfg.ServiceToken,
			Weight:       &r.cfg.Weight,
			Metadata:     rs.Metadata,
			Protocol:     &rs.Protocol,
			Version:      &r.cfg.Version,
		},
	}
	if !r.cfg.DisableHealthCheck {
		req.SetTTL(r.cfg.TTL)
	}
	resp, err := r.provider.Register(req)
	if err != nil {
		return fmt.Errorf("fail to Register instance, err is %v", err)
	}
	rs.InstanceID = resp.InstanceID
	return nil
}

func (r *Registry) heartBeat(rs *RegistryService) {
	tick := time.Second * time.Duration(r.cfg.HeartBeat)
	go func() {
		for {
			time.Sleep(tick)
			req := &api.InstanceHeartbeatRequest{
				InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
					Service:      rs.ServiceName,
					ServiceToken: r.cfg.ServiceToken,
					Namespace:    rs.Namespace,
					InstanceID:   rs.InstanceID,
					Host:         rs.Host,
					Port:         rs.Port,
				},
			}
			if err := r.provider.Heartbeat(req); err != nil {
				log.Println("heartbeat report err: %v\n", err)
			}
		}
	}()
}

// Deregister 反注册
func (r *Registry) Deregister(sn string) error {
	if !r.cfg.EnableRegister {
		return nil
	}
	rs := r.cfg.Services[sn]
	req := &api.InstanceDeRegisterRequest{
		InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
			Service:      rs.ServiceName,
			Namespace:    rs.Namespace,
			InstanceID:   rs.InstanceID,
			ServiceToken: r.cfg.ServiceToken,
			Host:         rs.Host,
			Port:         rs.Port,
		},
	}
	if err := r.provider.Deregister(req); err != nil {
		return fmt.Errorf("deregister error: %s", err.Error())
	}
	return nil
}
