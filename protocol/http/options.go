package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhengheng7913/grpc-config/server"
)

type Option struct {
	serviceName    string
	serveMuxOption runtime.ServeMuxOption
}

func (o Option) Apply(arg ...interface{}) {
	panic("implement me")
}

func (o Option) ServiceName() string {
	return o.serviceName
}

func (o Option) ProtocolName() string {
	return ProtocolName
}

func assertOptions(inters ...interface{}) ([]server.Option, bool) {
	opts := make([]server.Option, len(inters))
	for _, inter := range inters {
		opt, ok := inter.(server.Option)
		if !ok {
			return nil, false
		}
		opts = append(opts, opt)
	}
	return opts, true
}

func dessertOptions(opts ...runtime.ServeMuxOption) []interface{} {
	inters := make([]interface{}, len(opts))
	for _, opt := range opts {
		inters = append(inters, opt)
	}
	return inters
}
