package protocol

import (
	"github.com/zhengheng7913/grpc-config/protocol/grpc"
	"github.com/zhengheng7913/grpc-config/server"
)

var (
	ImplementMap = make(map[string]server.ServiceConstructor)
)

// Register 非线程安全
func Register(name string, constructor server.ServiceConstructor) {
	ImplementMap[name] = constructor
}

func Get(name string) server.ServiceConstructor {
	return ImplementMap[name]
}

func init() {
	Register(grpc.ProtocolName, grpc.NewGrpcService)
}
