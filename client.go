package grpc_go_starter

import (
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/client"
)

func NewClients(opt ...client.Option) *client.Clients {
	cfg, err := LoadSetup()
	if err != nil {
		panic(fmt.Errorf("load setup error: %s", err))
	}
	return NewClientWithConfig(cfg, opt...)
}

func newClientWithConfig(cfg *Config, conf *ServiceConfig, opt ...client.Option) *client.Clients {

}

func NewClientWithConfig(cfg *Config, opts ...client.Option) *client.Clients {

}
