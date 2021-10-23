package grpc

import (
	"github.com/zhengheng7913/grpc-config/server"
	"google.golang.org/grpc"
)

type Option struct {
	serviceName string
	opts        []grpc.ServerOption
}

func (o Option) ProtocolName() string {
	return ProtocolName
}

func (o Option) ServiceName() string {
	return o.serviceName
}

func (o *Option) Apply(inters ...interface{}) {
	gOpts, ok := assertOptions(inters)
	if !ok {
		panic("unknown service type")
	}
	o.opts = append(o.opts, gOpts...)
}

func assertOptions(inters ...interface{}) ([]grpc.ServerOption, bool) {
	opts := make([]grpc.ServerOption, len(inters))
	for _, inter := range inters {
		opt, ok := inter.(grpc.ServerOption)
		if !ok {
			return nil, false
		}
		opts = append(opts, opt)
	}
	return opts, true
}

func DessertOptions(opts ...grpc.ServerOption) []interface{} {
	inters := make([]interface{}, len(opts))
	for _, opt := range opts {
		inter := opt.(interface{})
		inters = append(inters, inter)
	}
	return inters
}

func WithGrpcOptions(serviceName string, option ...grpc.ServerOption) server.Options {
	return func(opt server.Option) {
		// 判断是本服务需要加载的option
		if serviceName != opt.ServiceName() {
			return
		}
		opt.Apply(DessertOptions(option...))
	}
}
